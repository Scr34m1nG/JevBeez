import { describe, expect, it } from 'bun:test'
import fs from 'node:fs'
import path from 'node:path'

describe('JevBeez Brand Logo & Assets', () => {
  const rootDir = path.resolve(import.meta.dir, '..')

  it('AppLogo.vue component exists with accessible SVG markup', () => {
    const logoPath = path.join(rootDir, 'components', 'AppLogo.vue')
    expect(fs.existsSync(logoPath)).toBe(true)

    const content = fs.readFileSync(logoPath, 'utf8')
    expect(content).toContain('<svg')
    expect(content).toContain('viewBox="0 0 64 64"')
    expect(content).toContain('role="img"')
    expect(content).toContain('aria-label="JevBeez Logo"')
    expect(content).toContain('jb-honey-gold')
    expect(content).toContain('jb-honey-bg')
  })

  it('favicon.svg exists in public directory and contains bee vector graphics', () => {
    const svgPath = path.join(rootDir, 'public', 'favicon.svg')
    expect(fs.existsSync(svgPath)).toBe(true)

    const content = fs.readFileSync(svgPath, 'utf8')
    expect(content).toContain('<svg')
    expect(content).toContain('xmlns="http://www.w3.org/2000/svg"')
    expect(content).toContain('jb-honey-gold')
    expect(content).toContain('jb-honey-bg')
  })

  it('favicon.ico exists and has valid ICO binary header', () => {
    const icoPath = path.join(rootDir, 'public', 'favicon.ico')
    expect(fs.existsSync(icoPath)).toBe(true)

    const buf = fs.readFileSync(icoPath)
    expect(buf.length).toBeGreaterThan(100)
    expect(buf.readUInt16LE(0)).toBe(0) // Reserved (0)
    expect(buf.readUInt16LE(2)).toBe(1) // Type (1 = ICO)
    expect(buf.readUInt16LE(4)).toBe(1) // Image count (1)
  })

  it('Navbar.vue integrates AppLogo', () => {
    const navbarPath = path.join(rootDir, 'components', 'Navbar.vue')
    const content = fs.readFileSync(navbarPath, 'utf8')
    expect(content).toContain('<AppLogo')
    expect(content).not.toContain('>JB</div>')
  })

  it('nuxt.config.ts configures SVG and ICO favicons in head', () => {
    const nuxtConfigPath = path.join(rootDir, 'nuxt.config.ts')
    const content = fs.readFileSync(nuxtConfigPath, 'utf8')
    expect(content).toContain("href: '/favicon.svg'")
    expect(content).toContain("href: '/favicon-32x32.png'")
    expect(content).toContain("href: '/favicon.ico'")
    expect(content).toContain("href: '/apple-touch-icon.png'")
  })

  it('favicon-32x32.png and apple-touch-icon.png exist in public directory', () => {
    const png32 = path.join(rootDir, 'public', 'favicon-32x32.png')
    const touchIcon = path.join(rootDir, 'public', 'apple-touch-icon.png')
    expect(fs.existsSync(png32)).toBe(true)
    expect(fs.existsSync(touchIcon)).toBe(true)
  })
})
