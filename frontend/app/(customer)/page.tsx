'use client'

import { useEffect } from "react"

export default function Page() {
  useEffect(() => {
    fetch('/api/booking/vehicle-types')
      .then(res => res.json())
      .then(data => console.log(data))
      .catch(err => console.log(err))
  }, [])
  
  return (
    <div>
      <h1>Customer App frontend</h1>
    </div>
  )
}
