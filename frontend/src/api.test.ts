import { afterEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './api'

afterEach(() => vi.unstubAllGlobals())

describe('calculate', () => {
  it('posts the exact square-root request and returns the result', async () => {
    const fetchMock = vi.fn(async (url: string, init: RequestInit) => {
      expect(url).toBe('/calculate')
      expect(init.method).toBe('POST')
      expect(init.headers).toEqual({ 'Content-Type': 'application/json' })
      expect(JSON.parse(String(init.body))).toEqual({ operation: 'sqrt', a: 9 })
      return Response.json({ result: 3 })
    })
    vi.stubGlobal('fetch', fetchMock)

    await expect(calculate({ operation: 'sqrt', a: 9 })).resolves.toBe(3)
    expect(fetchMock).toHaveBeenCalledOnce()
  })

  it('shows the backend error message', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => Response.json(
      { error: { code: 'division_by_zero', message: 'division by zero' } },
      { status: 400 },
    )))

    await expect(calculate({ operation: 'divide', a: 1, b: 0 })).rejects.toThrow('division by zero')
  })

  it('explains a network failure', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')))

    await expect(calculate({ operation: 'add', a: 1, b: 2 })).rejects.toThrow('Check that the backend is running')
  })
})
