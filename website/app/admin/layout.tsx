import type { Metadata } from 'next'
import '../globals.css'

export const metadata: Metadata = {
  title: 'JayPlus Customer App',
  description: 'Some Description',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang='en'>
      <body>
        <nav>This is Admin app of website</nav>
        {children}
      </body>
    </html>
  )
}
