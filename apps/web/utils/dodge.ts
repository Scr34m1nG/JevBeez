export interface DodgeOffset {
  x: number
  y: number
}

export function calculateDodgeOffset(dodgeCount: number): DodgeOffset {
  const directions = [
    { x: 35, y: -22 },
    { x: -20, y: 18 },
    { x: 45, y: 15 },
    { x: -15, y: -20 }
  ]
  return directions[(dodgeCount - 1) % directions.length]
}

export function getButtonTextForDodge(count: number): string {
  if (count === 0) return 'Ya, periksa!'
  if (count === 1) return 'Eits...!'
  if (count === 2) return 'Kejar aku!'
  return 'Hampir!'
}
