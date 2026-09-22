import { useState, type FormEvent } from 'react'
import { calculate } from './api'
import { operations, validateOperands, type Operation } from './calculator'

function App() {
  const [firstOperand, setFirstOperand] = useState('10')
  const [secondOperand, setSecondOperand] = useState('5')
  const [operation, setOperation] = useState<Operation>('add')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const isSquareRoot = operation === 'sqrt'

  function clearOutput() {
    setResult(null)
    setError('')
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const validation = validateOperands(firstOperand, secondOperand, operation)
    if (!validation.ok) {
      setResult(null)
      setError(validation.error)
      return
    }

    setError('')
    setResult(null)
    setIsLoading(true)

    try {
      setResult(await calculate(validation.request))
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : 'Something went wrong. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <main className="relative flex min-h-dvh items-center justify-center overflow-hidden bg-[#f0edfa] px-4 py-8 text-slate-900 sm:px-6 sm:py-12">
      <div aria-hidden="true" className="pointer-events-none absolute -top-32 -left-24 h-80 w-80 rounded-full bg-violet-300/45 blur-3xl" />
      <div aria-hidden="true" className="pointer-events-none absolute -right-28 -bottom-28 h-96 w-96 rounded-full bg-fuchsia-300/35 blur-3xl" />

      <section className="relative w-full max-w-lg overflow-hidden rounded-3xl bg-white shadow-[0_24px_80px_-20px_rgba(66,38,115,0.38)] ring-1 ring-violet-200/60" aria-labelledby="calculator-title">
        <header className="bg-gradient-to-br from-violet-700 via-purple-700 to-indigo-800 px-6 py-7 text-white sm:px-8">
          <div className="mb-3 flex items-center gap-2 text-xs font-semibold tracking-[0.22em] text-violet-200 uppercase">
            <span aria-hidden="true" className="inline-block h-2 w-2 rounded-full bg-lime-300 shadow-[0_0_12px_rgba(190,242,100,0.8)]" />
            Simple math
          </div>
          <h1 id="calculator-title" className="text-3xl font-bold tracking-tight sm:text-4xl">Calculator</h1>
          <p className="mt-2 text-sm text-violet-100">Choose an operation and get your result.</p>
        </header>

        <form onSubmit={handleSubmit} noValidate className="space-y-6 px-6 py-7 sm:px-8 sm:py-8">
          <div>
            <label htmlFor="first-operand" className="mb-2 block text-sm font-semibold text-slate-700">First operand</label>
            <input
              id="first-operand"
              name="firstOperand"
              type="number"
              step="any"
              inputMode="decimal"
              value={firstOperand}
              onChange={(event) => { setFirstOperand(event.target.value); clearOutput() }}
              disabled={isLoading}
              className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3.5 text-lg font-medium tabular-nums text-slate-900 outline-none transition focus:border-violet-500 focus:bg-white focus:ring-4 focus:ring-violet-100 disabled:cursor-wait disabled:opacity-60"
            />
          </div>

          <fieldset disabled={isLoading}>
            <legend className="mb-3 text-sm font-semibold text-slate-700">Operation</legend>
            <div className="grid grid-cols-4 gap-2.5 sm:gap-3">
              {operations.map(({ value, symbol, label }) => (
                <button
                  key={value}
                  type="button"
                  aria-label={label}
                  aria-pressed={operation === value}
                  onClick={() => { setOperation(value); clearOutput() }}
                  className={`flex h-14 items-center justify-center rounded-xl border text-xl font-semibold transition focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 disabled:cursor-wait ${operation === value ? 'border-violet-700 bg-violet-700 text-white shadow-lg shadow-violet-200' : 'border-slate-200 bg-white text-slate-700 hover:border-violet-300 hover:bg-violet-50'}`}
                >
                  {value === 'power' ? <span aria-hidden="true">x<sup className="text-sm">y</sup></span> : <span aria-hidden="true">{symbol}</span>}
                </button>
              ))}
            </div>
            {operation === 'percentage' && <p className="mt-2 text-xs text-slate-500">First operand % of second operand.</p>}
          </fieldset>

          <div>
            <label htmlFor="second-operand" className="mb-2 block text-sm font-semibold text-slate-700">Second operand</label>
            <input
              id="second-operand"
              name="secondOperand"
              type="number"
              step="any"
              inputMode="decimal"
              value={secondOperand}
              onChange={(event) => { setSecondOperand(event.target.value); clearOutput() }}
              disabled={isSquareRoot || isLoading}
              aria-describedby={isSquareRoot ? 'second-operand-hint' : undefined}
              className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3.5 text-lg font-medium tabular-nums text-slate-900 outline-none transition focus:border-violet-500 focus:bg-white focus:ring-4 focus:ring-violet-100 disabled:cursor-not-allowed disabled:bg-slate-100 disabled:text-slate-400"
            />
            {isSquareRoot && <p id="second-operand-hint" className="mt-2 text-xs text-slate-500">Square root uses only the first operand.</p>}
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="flex w-full items-center justify-center rounded-xl bg-violet-700 px-5 py-3.5 text-base font-semibold text-white shadow-lg shadow-violet-200 transition hover:bg-violet-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 disabled:cursor-wait disabled:opacity-70"
          >
            {isLoading ? 'Calculating…' : 'Calculate'}
          </button>

          {error && <p role="alert" className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">{error}</p>}

          <div className="rounded-2xl border border-violet-100 bg-violet-50 px-5 py-5" aria-live="polite" aria-busy={isLoading}>
            <h2 className="text-xs font-bold tracking-[0.16em] text-violet-700 uppercase">Result</h2>
            <output className="mt-2 block min-h-10 break-all text-3xl font-bold tabular-nums text-slate-900 sm:text-4xl">
              {result === null ? <span className="text-slate-400">—</span> : String(result)}
            </output>
          </div>
        </form>
      </section>
    </main>
  )
}

export default App
