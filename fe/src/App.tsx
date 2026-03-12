import { useState } from 'react'
import { Button, Container, Typography } from '@mui/material'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-slate-50 py-10">
      <Container maxWidth="sm">
        <Typography variant="h4" component="h1" gutterBottom>
          Kiosk
        </Typography>

        <div className="rounded-lg bg-white p-6 shadow">
          <Typography variant="body1" gutterBottom>
            Count: {count}
          </Typography>
          <Button variant="contained" onClick={() => setCount((c) => c + 1)}>
            Increment
          </Button>
        </div>
      </Container>
    </div>
  )
}

export default App
