import { describe, expect, it } from 'vitest'
import { validateOperands } from './calculator'

describe('validateOperands', () => {
  it('keeps zero as a valid operand', () => {
    expect(validateOperands('0', '5', 'subtract')).toEqual({
      ok: true,
      request: { operation: 'subtract', a: 0, b: 5 },
    })
  })

  it('omits the second operand for square root even when it has a value', () => {
    expect(validateOperands('9', '5', 'sqrt')).toEqual({
      ok: true,
      request: { operation: 'sqrt', a: 9 },
    })
  })

  it('requires both operands for binary operations', () => {
    expect(validateOperands('', '5', 'add')).toEqual({ ok: false, error: 'First operand is required.' })
    expect(validateOperands('5', '', 'add')).toEqual({ ok: false, error: 'Second operand is required.' })
  })

  it('rejects non-finite and non-decimal values', () => {
    expect(validateOperands('1e309', '5', 'multiply').ok).toBe(false)
    expect(validateOperands('0x10', '5', 'multiply').ok).toBe(false)
  })
})
