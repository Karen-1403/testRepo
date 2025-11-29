export function ZGateLogo({ size = 48, className = "" }: { size?: number; className?: string }) {
  return (
    <div
      className={`relative overflow-hidden flex-shrink-0 rounded-full flex items-center justify-center shadow-lg shadow-cyan-500/50 ${className}`}
      style={{
        width: size,
        height: size,
      }}
    >
      <img src="/images/image.png" alt="zGate Logo" className="w-full h-full object-cover" />
    </div>
  )
}
