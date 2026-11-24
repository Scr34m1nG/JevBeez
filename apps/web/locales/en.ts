import type { TranslationSchema } from './id'

export const en: TranslationSchema = {
  nav: {
    brandBadge: 'Fact Check',
    systemActive: 'System Active',
    engineName: 'TypeSafe Jev'
  },
  hero: {
    badge: 'Automated AI Fact-Checker',
    title: 'Verify Claims & News Facts',
    description: 'Paste a news link from chats or social media to verify claims with evidence-grounded AI decision modeling.'
  },
  form: {
    inputLabel: 'Paste a news or social media link:',
    placeholder: 'https://...',
    pasteTitle: 'Paste from Clipboard',
    pasteBtn: 'Paste',
    clearAria: 'Clear link',
    checkBtn: 'Fact-Check',
    checkingBtn: 'Checking Facts...',
    blockedBtn: 'Access Blocked',
    cooldownWaiting: 'Queued ({time}s)...',
    cooldownWaitBtn: 'Wait a moment ({time}s)',
    cooldownNotice: 'Please wait {time} seconds before the next check (or click the button to automatically check after the cooldown).',
    blockedTitle: 'Access Denied: Your IP is Blocked (24 Hours)',
    blockedReasonLabel: 'Reason:',
    easterEggTitle: 'Whoops, Mirror Inception! 🤖',
    easterEggBtn: 'Got it, check another site 🫡',
    queueAlertTitle: 'Link is in queue:',
    queueAlertDesc: 'Will automatically check in {time} seconds...',
    cancelQueue: 'Cancel'
  },
  results: {
    checkAnotherBtn: 'Check Another Link',
    cachedResult: 'Cached Result',
    confidenceLabel: 'Confidence Level:',
    sourcePrefix: 'Source:',
    socialMedia: 'Social Media',
    summaryTitle: 'Fact Analysis Summary:',
    shareWhatsApp: 'Share to WhatsApp',
    copyBtn: 'Copy Summary',
    copiedBtn: 'Copied!',
    redFlagsTitle: 'Suspicious Red Flags',
    detailedHeading: 'AI Decision Details',
    contradictedByEvidence: 'Contradicted by Evidence',
    supportedByEvidence: 'Supported by Evidence',
    sensationalFraming: 'Sensational / Clickbait Framing',
    shareMsgTitle: '📢 *JevBeez AI Fact-Check Result*',
    shareMsgStatus: '*Status:*',
    shareMsgConfidence: '*Confidence Level:*',
    shareMsgArticleTitle: '*Headline:*',
    shareMsgSummary: '*Explanation:*',
    shareMsgCheckAt: '*Check truth of news at:*',
    copyHeader: '[JevBeez Fact-Check]',
    copyStatus: 'Status:',
    copyConfidence: 'AI Confidence:',
    copyTitle: 'Title:',
    copySummaryLabel: 'Summary:',
    copyUrlLabel: 'Article URL:'
  },
  verdict: {
    hoax: 'Confirmed Hoax / Fake',
    valid: 'Verified Fact',
    misleading: 'Misleading Content',
    unverified: 'Unverified Info',
    phishing: 'DANGER: PHISHING DETECTED',
    malware: 'DANGER: MALWARE / APK DETECTED',
    scam: 'WARNING: SCAM DETECTED'
  },
  threat: {
    threatBadge: 'Cyber Threat Detected',
    anomalyTitle: 'Link & Domain Structural Anomalies',
    actionAdviceTitle: 'Recommended Security Actions',
    action1: 'Never enter your password, card number, PIN, or OTP codes on this link.',
    action2: 'If you already downloaded an APK or executable file, do not open it and delete it immediately.',
    action3: 'If impersonating a bank or authority, verify directly via their official customer service channels.',
    deepDiveBtn: 'Run In-Depth Content Audit (Sandboxed Deep Dive)',
    deepDiveRunning: 'Inspecting forms and content in sandbox...',
    deepDiveCompleted: 'Audited via Sandboxed Deep Dive',
    formAlertTitle: 'Credential Harvesting Form Indicators',
    unshortenedNotice: 'Original URL unrolled from shortener',
    targetBrandPrefix: 'Impersonating Brand:'
  },
  evidence: {
    title: 'Sources & Evidence',
    badge: 'Verified Web Search',
    openSource: 'Open Source'
  },
  selfCheck: {
    modalTitle: 'Still want to try?',
    modalDesc: 'You are checking JevBeez on JevBeez itself. Are you sure you want to continue?',
    cancelBtn: 'Never mind',
    dodge0: 'Yes, check!',
    dodge1: 'Whoops...!',
    dodge2: 'Catch me!',
    dodge3: 'Almost!',
    taunt1: 'Whoops! Too fast for you 💨',
    taunt2: 'Haha 😜 Too slow on the click!',
    taunt3: 'Almost got it! Try again 🏃‍♂️',
    taunt4: 'Nimble button, right? 🎯'
  },
  footer: {
    liveVerification: 'JevBeez • Automated real-time evidence-grounded fact checking.',
    techStack: 'Nuxt 3 & Go Fiber Core • Powered by TypeSafe Jev'
  }
}
