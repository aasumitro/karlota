package common

const ForgotPasswordTemplate = `
<div style="margin-bottom: 8px;">
<p style="margin-bottom:2px; font-weight:bold;">Hi %[1]s,</p>
<p>You have received this email because a password reset request for your account was received.</p>
</div>

<div style="display: flex; justify-content: center; align-items: center; margin-bottom: 16px;">
  <a class="btn" href="%[2]s">Reset Password</a>
</div>

<hr />
<div style="margin-top: 8px;"><i>If you’re having trouble with the button 'Reset Password', copy and paste (or click) the URL below into your web browser.</i></div>
<a href="%[2]s" style="font-size: 10px; color: blue;">%[2]s</a>
`

const VerifyEmailTemplate = `
<div style="margin-bottom: 8px;">
<p style="margin-bottom:2px; font-weight:bold;">Hi %[1]s,</p>
<p>Verify your email address to complete your registration.</p>
</div>

<div style="display: flex; justify-content: center; align-items: center; margin-bottom: 16px;">
  <a class="btn" href="%[2]s">Verify Email</a>
</div>

<hr />
<div style="margin-top: 8px;"><i>If you’re having trouble with the button 'Verify Email', copy and paste (or click) the URL below into your web browser.</i></div>
<a href="%[2]s" style="font-size: 10px; color: blue;">%[2]s</a>
`
