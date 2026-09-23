import type { CalculationRequest } from './calculator'

type ApiError = { error: { code: string; message: string } }
type ApiSuccess = { result: number }

function isApiError(value: unknown): value is ApiError {
  return typeof value === 'object' && value !== null &&
    'error' in value && typeof value.error === 'object' && value.error !== null &&
    'message' in value.error && typeof value.error.message === 'string'
}

function isApiSuccess(value: unknown): value is ApiSuccess {
  return typeof value === 'object' && value !== null &&
    'result' in value && typeof value.result === 'number' && Number.isFinite(value.result)
}

export async function calculate(request: CalculationRequest): Promise<number> {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 15_000)

  try {
    let response: Response
    try {
      response = await fetch('/calculate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
        signal: controller.signal,
      })
    } catch {
      if (controller.signal.aborted) throw new Error('The calculation timed out. Please try again.')
      throw new Error('Could not reach the calculator. Check that the backend is running.')
    }

    let data: unknown = null
    try {
      data = await response.json()
    } catch {
      if (controller.signal.aborted) throw new Error('The calculation timed out. Please try again.')
    }

    if (!response.ok) {
      throw new Error(isApiError(data) ? data.error.message : 'Calculation failed. Please try again.')
    }

    if (!isApiSuccess(data)) throw new Error('The calculator returned an unexpected response.')

    return data.result
  } finally {
    clearTimeout(timeout)
  }
}
