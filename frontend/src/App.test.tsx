// @vitest-environment jsdom

import { act } from 'react'
import { createRoot, type Root } from 'react-dom/client'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'

Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true })

let container: HTMLDivElement
let root: Root

beforeEach(async () => {
  container = document.createElement('div')
  document.body.append(container)
  root = createRoot(container)
  await act(async () => root.render(<App />))
})

afterEach(async () => {
  await act(async () => root.unmount())
  container.remove()
  vi.unstubAllGlobals()
})

async function submitForm() {
  const form = container.querySelector('form')
  if (!form) throw new Error('Calculator form was not rendered')
  await act(async () => {
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }))
  })
}

describe('Calculator UI', () => {
  it('keeps the second operand visible but disabled for square root, and omits it from the request', async () => {
    const fetchMock = vi.fn(async (url: string, init: RequestInit) => {
      expect(url).toBe('/calculate')
      expect(JSON.parse(String(init.body))).toEqual({ operation: 'sqrt', a: 10 })
      return Response.json({ result: Math.sqrt(10) })
    })
    vi.stubGlobal('fetch', fetchMock)

    const squareRootButton = container.querySelector<HTMLButtonElement>('button[aria-label="Square root"]')
    if (!squareRootButton) throw new Error('Square root button was not rendered')
    await act(async () => squareRootButton.click())

    const secondOperand = container.querySelector<HTMLInputElement>('#second-operand')
    expect(secondOperand).not.toBeNull()
    expect(secondOperand?.disabled).toBe(true)

    await submitForm()

    expect(fetchMock).toHaveBeenCalledOnce()
    expect(container.querySelector('output')?.textContent).toBe(String(Math.sqrt(10)))
  })

  it('displays an error returned by the backend', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => Response.json(
      { error: { code: 'division_by_zero', message: 'division by zero' } },
      { status: 400 },
    )))

    await submitForm()

    expect(container.querySelector('[role="alert"]')?.textContent).toBe('division by zero')
    expect(container.querySelector('output')?.textContent).toBe('—')
  })

  it('displays a validation error without calling the backend', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)

    const firstOperand = container.querySelector<HTMLInputElement>('#first-operand')
    if (!firstOperand) throw new Error('First operand was not rendered')
    const setValue = Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, 'value')?.set
    if (!setValue) throw new Error('Input value setter is unavailable')
    await act(async () => {
      setValue.call(firstOperand, '')
      firstOperand.dispatchEvent(new Event('input', { bubbles: true }))
    })

    await submitForm()

    expect(container.querySelector('[role="alert"]')?.textContent).toBe('First operand is required.')
    expect(fetchMock).not.toHaveBeenCalled()
  })
})
