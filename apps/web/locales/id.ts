export const id = {
  nav: {
    brandBadge: 'Cek Fakta',
    systemActive: 'Sistem Aktif',
    engineName: 'TypeSafe Jev'
  },
  hero: {
    badge: 'Pemeriksa Fakta Otomatis AI',
    title: 'Cek Kebenaran Fakta & Berita',
    description: 'Tempel link berita dari grup WA atau medsos buat buktiin beneran fakta atau cuma hoax berbasis bukti dan model keputusan Jev.'
  },
  form: {
    inputLabel: 'Tempel link berita atau medsos:',
    placeholder: 'https://...',
    pasteTitle: 'Tempel dari Clipboard',
    pasteBtn: 'Tempel',
    clearAria: 'Hapus tautan',
    checkBtn: 'Cek Fakta',
    checkingBtn: 'Memeriksa Fakta...',
    blockedBtn: 'Akses Diblokir',
    cooldownWaiting: 'Lagi antre ({time}s)...',
    cooldownWaitBtn: 'Tunggu sebentar ({time}s)',
    cooldownNotice: 'Mohon tunggu jeda {time} detik sebelum pemeriksaan berikutnya (atau klik tombol untuk proses otomatis setelah jeda).',
    blockedTitle: 'Akses Ditolak: Alamat IP Anda Diblokir (24 Jam)',
    blockedReasonLabel: 'Alasan:',
    easterEggTitle: 'Waduh, Kena Prank Cermin! 🤖',
    easterEggBtn: 'Siap, coba web lain 🫡',
    queueAlertTitle: 'Link kamu udah masuk antrean:',
    queueAlertDesc: 'Bakal otomatis dicek dalam {time} detik...',
    cancelQueue: 'Batalkan',
    modeUrl: 'Mode Tautan',
    modeText: 'Mode Teks Pesan',
    modeAuto: 'Deteksi Otomatis'
  },
  results: {
    checkAnotherBtn: 'Cek Link Lain',
    cachedResult: 'Hasil Tersimpan (Cache)',
    confidenceLabel: 'Tingkat Keyakinan:',
    sourcePrefix: 'Sumber:',
    socialMedia: 'Media Sosial',
    summaryTitle: 'Ringkasan Analisis Fakta:',
    shareWhatsApp: 'Bagikan ke WhatsApp',
    copyBtn: 'Salin Ringkasan',
    copiedBtn: 'Udah Disalin!',
    copyWaBtn: 'Salin Format WA',
    copiedWaBtn: 'Format WA Disalin!',
    downloadCardBtn: 'Download Kartu Bukti',
    downloadingCardBtn: 'Menyiapkan Kartu...',
    redFlagsTitle: 'Hal yang Mencurigakan',
    detailedHeading: 'Detail Pengujian Fakta AI',
    contradictedByEvidence: 'Indikasi Bohong / Bantahan',
    supportedByEvidence: 'Didukung Bukti Valid',
    sensationalFraming: 'Gaya Bahasa Klikbait / Sensasi',
    shareMsgTitle: '📢 *Hasil Cek Fakta JevBeez AI*',
    shareMsgStatus: '*Status:*',
    shareMsgConfidence: '*Tingkat Keyakinan:*',
    shareMsgArticleTitle: '*Judul Berita:*',
    shareMsgSummary: '*Penjelasan:*',
    shareMsgCheckAt: '*Cek kebenaran berita di:*',
    copyHeader: '[Pemeriksaan Fakta JevBeez]',
    copyStatus: 'Status:',
    copyConfidence: 'Keyakinan AI:',
    copyTitle: 'Judul:',
    copySummaryLabel: 'Ringkasan:',
    copyUrlLabel: 'Tautan Berita:'
  },
  verdict: {
    hoax: 'Terbukti Hoax / Palsu',
    valid: 'Beneran Fakta (Akurat)',
    misleading: 'Menyesatkan / Misleading',
    unverified: 'Belum Bisa Dipastikan',
    phishing: 'BAHAYA: INDIKASI PHISHING',
    malware: 'BAHAYA: UNDUHAN MALWARE / APK',
    scam: 'WASPADA: SITUS PENIPUAN (SCAM)',
    ads: 'Iklan Berbayar (Ads)'
  },
  threat: {
    threatBadge: 'Ancaman Siber Terdeteksi',
    anomalyTitle: 'Keanehan Struktur Tautan & Domain',
    actionAdviceTitle: 'Rekomendasi Tindakan Keamanan',
    action1: 'Jangan masukkan password, nomor kartu, PIN, atau kode OTP di tautan ini.',
    action2: 'Jika telah terlanjur mengunduh file APK/eksekusi, jangan buka dan segera hapus dari folder Download.',
    action3: 'Jika mengatasnamakan bank atau instansi, hubungi saluran pengaduan resmi institusi terkait.',
    deepDiveBtn: 'Buka Analisis Konten Mendalam (Sandboxed Deep Dive)',
    deepDiveRunning: 'Sedang menganalisis konten dalam sandbox...',
    deepDiveCompleted: 'Hasil Audit Deep Dive Sandboxed',
    formAlertTitle: 'Indikator Formulir Pencurian Data',
    unshortenedNotice: 'Tautan asli telah dibongkar dari pemendek URL',
    targetBrandPrefix: 'Mencatut Brand:'
  },
  evidence: {
    title: 'Sumber & Bukti Terkait',
    badge: 'Penelusuran Web Terverifikasi',
    openSource: 'Buka Sumber'
  },
  selfCheck: {
    modalTitle: 'Yakin masih mau coba?',
    modalDesc: 'Anda sudah mencoba mengecek website JevBeez di JevBeez sendiri. Apakah Anda yakin ingin melanjutkan?',
    cancelBtn: 'Hmm, nggak jadi deh',
    dodge0: 'Ya, periksa!',
    dodge1: 'Eits...!',
    dodge2: 'Kejar aku!',
    dodge3: 'Hampir!',
    taunt1: 'Eits! Kecepetan tombolnya kabur 💨',
    taunt2: 'Wlee 😜 Kurang cepet ngeklik!',
    taunt3: 'Hampir aja kena! Coba kejar lagi 🏃‍♂️',
    taunt4: 'Lincah kan tombolnya? 🎯'
  },
  footer: {
    liveVerification: 'JevBeez • Verifikasi fakta otomatis berbasis bukti live.',
    techStack: 'Nuxt 3 & Go Fiber Core • Powered by TypeSafe Jev'
  }
}

export type TranslationSchema = typeof id
