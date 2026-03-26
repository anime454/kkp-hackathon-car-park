import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import './App.css'

type SlotStatus = 'free' | 'parked' | 'close'
type SlotType = 'VIP' | 'normal'

type Slot = {
  id: string
  label: string
  type: SlotType
  status: SlotStatus
  parkedAt: string | null
}

type Dashboard = {
  slots: Slot[]
  total: number
  free: number
  parked: number
  closed: number
  updatedAt: string
  priceNormal: number
  priceVip: number
}

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

async function fetchDashboard(): Promise<Dashboard> {
  const res = await fetch(`${API_BASE}/kiosk/dashboard`)
  if (!res.ok) throw new Error('Failed to fetch dashboard')
  return res.json()
}

async function fetchKioskInfo() {
  const res = await fetch(`${API_BASE}/kiosk/info`)
  if (!res.ok) throw new Error('Failed to fetch kiosk info')
  return res.json() as Promise<{
    message: string
    vipRate: string
    normalRate: string
    availableSlot: number
    occupiedSlot: number
    updatedAt: string
  }>
}

async function adminLogin(username: string, password: string) {
  const res = await fetch(`${API_BASE}/management/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  if (!res.ok) throw new Error('Invalid credentials')
  return res.json() as Promise<{ token: string }>
}

async function fetchManagementSlots(token: string) {
  const res = await fetch(`${API_BASE}/management/slots`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) throw new Error('Unauthorized or failed to fetch slots')
  return res.json() as Promise<{ slots: Slot[] }>
}

async function updateSlotStatus(token: string, slotId: string, status: SlotStatus) {
  const res = await fetch(`${API_BASE}/management/slots/${slotId}/status`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ status }),
  })
  const body = await res.json().catch(() => null)
  if (!res.ok) {
    throw new Error((body && body.error) || 'Failed to update slot')
  }
  return body as { fee?: { amount: number; hours: number }; slot: Slot }
}

function App() {
  const queryClient = useQueryClient()
  const [mode, setMode] = useState<'kiosk' | 'management'>('kiosk')
  const [token, setToken] = useState<string>('')
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin123')
  const [message, setMessage] = useState('')

  const dashboardQuery = useQuery({
    queryKey: ['dashboard'],
    queryFn: fetchDashboard,
    refetchInterval: 3000,
  })

  const kioskInfoQuery = useQuery({
    queryKey: ['kiosk-info'],
    queryFn: fetchKioskInfo,
    refetchInterval: 7000,
  })

  const managementSlotsQuery = useQuery({
    queryKey: ['management-slots', token],
    queryFn: () => fetchManagementSlots(token),
    enabled: Boolean(token),
    refetchInterval: 3000,
  })

  const loginMutation = useMutation({
    mutationFn: () => adminLogin(username, password),
    onSuccess: (data) => {
      setToken(data.token)
      setMessage('Logged in as admin')
    },
    onError: (error) => {
      setMessage(error instanceof Error ? error.message : 'Login failed')
    },
  })

  const updateStatusMutation = useMutation({
    mutationFn: ({ slotId, status }: { slotId: string; status: SlotStatus }) =>
      updateSlotStatus(token, slotId, status),
    onSuccess: (data) => {
      if (data.fee) {
        setMessage(`Parking fee: ${data.fee.amount} (${data.fee.hours} hour(s))`)
      } else {
        setMessage('Slot updated')
      }
      queryClient.invalidateQueries({ queryKey: ['management-slots'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    },
    onError: (error) => {
      setMessage(error instanceof Error ? error.message : 'Update failed')
    },
  })

  const slotStats = useMemo(() => {
    const data = dashboardQuery.data
    if (!data) return { freePct: 0, occupiedPct: 0 }
    const total = data.total || 1
    return {
      freePct: Math.round((data.free / total) * 100),
      occupiedPct: Math.round((data.parked / total) * 100),
    }
  }, [dashboardQuery.data])

  return (
    <main className="app-shell">
      <section className="topbar">
        <h1>KKP Smart Car Park</h1>
        <div className="mode-switch">
          <button
            className={mode === 'kiosk' ? 'active' : ''}
            onClick={() => setMode('kiosk')}
          >
            Kiosk Mode
          </button>
          <button
            className={mode === 'management' ? 'active' : ''}
            onClick={() => setMode('management')}
          >
            Management Mode
          </button>
        </div>
      </section>

      {mode === 'kiosk' && (
        <section className="panel">
          <div className="kpi-grid">
            <article>
              <p>Free Slot</p>
              <h2>{dashboardQuery.data?.free ?? '-'}</h2>
            </article>
            <article>
              <p>Using Slot</p>
              <h2>{dashboardQuery.data?.parked ?? '-'}</h2>
            </article>
            <article>
              <p>Closed Slot</p>
              <h2>{dashboardQuery.data?.closed ?? '-'}</h2>
            </article>
            <article>
              <p>Availability</p>
              <h2>{slotStats.freePct}%</h2>
            </article>
          </div>

          <div className="rates">
            <div>
              <strong>Normal:</strong> {kioskInfoQuery.data?.normalRate ?? '10 first hour, +5 each additional hour'}
            </div>
            <div>
              <strong>VIP:</strong> {kioskInfoQuery.data?.vipRate ?? '20 first hour, +10 each additional hour'}
            </div>
            <div className="hint">{kioskInfoQuery.data?.message ?? 'Welcome to KKP Smart Parking'}</div>
          </div>

          <div className="slot-grid">
            {(dashboardQuery.data?.slots ?? []).map((slot) => (
              <article key={slot.id} className={`slot-card status-${slot.status}`}>
                <div className="slot-head">
                  <h3>{slot.label}</h3>
                  <span>{slot.type}</span>
                </div>
                <p>{slot.status.toUpperCase()}</p>
              </article>
            ))}
          </div>
        </section>
      )}

      {mode === 'management' && (
        <section className="panel">
          {!token ? (
            <form
              className="login-form"
              onSubmit={(e) => {
                e.preventDefault()
                loginMutation.mutate()
              }}
            >
              <h2>Admin Login</h2>
              <input value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" />
              <input
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                type="password"
                placeholder="Password"
              />
              <button type="submit" disabled={loginMutation.isPending}>
                {loginMutation.isPending ? 'Signing in...' : 'Sign in'}
              </button>
            </form>
          ) : (
            <>
              <div className="management-header">
                <h2>Slot Status Control</h2>
                <button
                  onClick={() => {
                    setToken('')
                    setMessage('Logged out')
                  }}
                >
                  Logout
                </button>
              </div>

              <div className="slot-grid">
                {(managementSlotsQuery.data?.slots ?? []).map((slot) => (
                  <article key={slot.id} className={`slot-card status-${slot.status}`}>
                    <div className="slot-head">
                      <h3>{slot.label}</h3>
                      <span>{slot.type}</span>
                    </div>
                    <p>{slot.status.toUpperCase()}</p>
                    <div className="actions">
                      <button onClick={() => updateStatusMutation.mutate({ slotId: slot.id, status: 'free' })}>
                        Free
                      </button>
                      <button onClick={() => updateStatusMutation.mutate({ slotId: slot.id, status: 'parked' })}>
                        Parked
                      </button>
                      <button onClick={() => updateStatusMutation.mutate({ slotId: slot.id, status: 'close' })}>
                        Close
                      </button>
                    </div>
                  </article>
                ))}
              </div>
            </>
          )}
        </section>
      )}

      <footer className="footer-note">
        <span>{message}</span>
        <span>
          Live occupancy: {slotStats.occupiedPct}% | Updated:{' '}
          {dashboardQuery.data ? new Date(dashboardQuery.data.updatedAt).toLocaleTimeString() : '-'}
        </span>
      </footer>
    </main>
  )
}

export default App
