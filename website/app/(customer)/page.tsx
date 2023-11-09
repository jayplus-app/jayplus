'use client'
import { use, useEffect } from "react"

export default function Page() {
  useEffect(() => {
    console.log('Customer page loaded in website')
  }, [])

  useEffect(() => {
    fetch('/api/hello')
      .then(response => response.json())
      .then(data => console.log(data))
      .catch(err => console.log(err))
  }, [])

  useEffect(() => {
    fetch('/api/hello/1')
      .then(response => response.json())
      .then(data => console.log(data))
      .catch(err => console.log(err))
  }, [])

  return (
    <div>
      <h1>Customer App</h1>
    </div>
  )
}
