export const operations = [
  { value: 'add', symbol: '+', label: 'Add' },
  { value: 'subtract', symbol: '−', label: 'Subtract' },
  { value: 'multiply', symbol: '×', label: 'Multiply' },
  { value: 'divide', symbol: '÷', label: 'Divide' },
  { value: 'power', symbol: '', label: 'Power' },
  { value: 'sqrt', symbol: '√', label: 'Square root' },
  { value: 'percentage', symbol: '%', label: 'Percentage' },
] as const

export type Operation = (typeof operations)[number]['value']

export type CalculationRequest = {
  operation: Operation
  a: number
  b?: number
}

type ValidationResult =
  | { ok: true; request: CalculationRequest }
  | { ok: false; error: string }

function parseOperand(raw: string, label: string): number | string {
  const trimmed = raw.trim()
  if (trimmed === '') return `${label} is required.`

  const decimal = /^[+-]?(?:\d+\.?\d*|\.\d+)(?:[eE][+-]?\d+)?$/
  const value = Number(trimmed)
  if (!decimal.test(trimmed) || !Number.isFinite(value)) return `${label} must be a finite number.`

  return value
}

export function validateOperands(first: string, second: string, operation: Operation): ValidationResult {
  const a = parseOperand(first, 'First operand')
  if (typeof a === 'string') return { ok: false, error: a }

  if (operation === 'sqrt') return { ok: true, request: { operation, a } }

  const b = parseOperand(second, 'Second operand')
  if (typeof b === 'string') return { ok: false, error: b }

  return { ok: true, request: { operation, a, b } }
}
