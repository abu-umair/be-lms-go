package utils

import "fmt"

func GetOTPEmailTemplate(code string) string {
	otpHTML := ""

	for _, d := range code {
		otpHTML += fmt.Sprintf(`<div class="otp-digit">%c</div>`, d)
	}

	return fmt.Sprintf(`
	<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Verifikasi Akun - MyApps</title>
  <style>
    body { margin: 0; padding: 20px; background: #f5f7fa; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }
    .container { max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1); }
    .header { background: #2563eb; padding: 24px 32px; text-align: center; }
    .header-content { display: flex; align-items: center; justify-content: center; gap: 12px; }
    .header-icon { width: 40px; height: 40px; background: rgba(255,255,255,0.2); border-radius: 8px; display: flex; align-items: center; justify-content: center; }
    .header h1 { color: #ffffff; font-size: 20px; font-weight: 600; margin: 0; letter-spacing: -0.02em; }
    .content { padding: 40px 32px; }
    .content p { color: #64748b; font-size: 16px; line-height: 1.6; margin: 0 0 32px; }
    .content p:first-child { color: #1e293b; margin-bottom: 8px; }
    .otp-box { background: #f1f5f9; border-radius: 12px; padding: 24px; margin-bottom: 32px; text-align: center; }
    .otp-label { font-size: 12px; font-weight: 500; color: #64748b; text-transform: uppercase; letter-spacing: 0.1em; margin-bottom: 12px; }
    .otp-digits { display: flex; justify-content: center; gap: 8px; }
    .otp-digit { width: 48px; height: 56px; background: #ffffff; border: 2px solid rgba(37,99,235,0.2); border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 24px; font-weight: 700; color: #2563eb; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); }
    .expiry-notice { display: flex; align-items: flex-start; gap: 12px; background: #fef3c7; border: 1px solid #fcd34d; border-radius: 8px; padding: 16px; margin-bottom: 24px; }
    .expiry-notice svg { flex-shrink: 0; color: #d97706; }
    .expiry-notice p { margin: 0; font-size: 14px; color: #92400e; }
    .expiry-notice strong { font-weight: 600; }
    .security-text { font-size: 14px !important; margin-bottom: 0 !important; }
    .footer { background: #f8fafc; padding: 24px 32px; border-top: 1px solid #e2e8f0; text-align: center; }
    .footer-company { font-size: 14px; font-weight: 500; color: #1e293b; margin: 0 0 4px; }
    .footer-auto { font-size: 12px; color: #64748b; margin: 0 0 16px; }
    .footer-links { display: flex; justify-content: center; gap: 16px; margin-bottom: 16px; }
    .footer-links a { font-size: 12px; color: #64748b; text-decoration: none; }
    .footer-links a:hover { color: #2563eb; }
    .footer-links span { color: #cbd5e1; }
    .copyright { font-size: 12px; color: #94a3b8; margin: 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="24" height="24" fill="none" stroke="white" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
          </svg>
        </div>
        <h1>Verifikasi Akun</h1>
      </div>
    </div>
    <div class="content">
      <p>Halo,</p>
      <p>Gunakan kode OTP di bawah ini untuk menyelesaikan proses verifikasi akun Anda. Jangan bagikan kode ini kepada siapapun.</p>
      <div class="otp-box">
        <div class="otp-label">Kode Verifikasi</div>
        <div class="otp-digits">
          %s
        </div>
      </div>
      <div class="expiry-notice">
        <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <p>Kode ini akan <strong>kadaluarsa dalam 5 menit</strong>. Segera masukkan kode untuk melanjutkan proses verifikasi.</p>
      </div>
      <p class="security-text">Jika Anda tidak merasa meminta kode ini, abaikan email ini atau hubungi tim support kami jika Anda mencurigai aktivitas yang tidak sah pada akun Anda.</p>
    </div>
    <div class="footer">
      <p class="footer-company">MyApps</p>
      <p class="footer-auto">Email ini dikirim secara otomatis, mohon tidak membalas email ini.</p>
      <div class="footer-links">
        <a href="#">Kebijakan Privasi</a>
        <span>•</span>
        <a href="#">Syarat & Ketentuan</a>
        <span>•</span>
        <a href="#">Bantuan</a>
      </div>
      <p class="copyright">© 2026 MyApps. All rights reserved.</p>
    </div>
  </div>
</body>
</html>
	`, otpHTML)
}
