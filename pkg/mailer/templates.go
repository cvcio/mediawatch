package mailer

var (
	TemplateMagicLink = `
<div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
	<p>
		Hi there!
	</p>
	<p>
		You requested an email link to sign into your MediaWatch account.<br />
		<a href="{{.baseURL}}/{{.link}}" target="_blank">Here is it!</a>
	</p>
	<p>
		The Team at MediaWatch
	</p>
	<p style="font-style:italic;">
		ps. If you need anything, we can help. And we also want your feedback – <br />
		good and bad, we want it all :) Get in touch at feedback@mediawatch.io.
	</p>
</div>
`
)

// package mailer

// var (
// 	msgInvitation = `
// <div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
//     <p>
//         Hi there
//     </p>
//     <p>
//         {{.name}} ({{.team}}) has invited you to join <a href="https://mediawatch.io" target="_blank">MediaWatch</a>.
//     </p>
//     <p>
//         If you are interested you can create an account by following this URL:<br/>
//         <a href="https://app.mediawatch.io/auth/register/{{.email}}/{{.nonce}}" target="_blank">https://app.mediawatch.io/auth/register/{{.email}}/{{.nonce}}</a>
//     </p>
//     <p>
//         If you are are not interested just ignore this message.<br/>
//         However, if you have any security or privacy discrimination concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
//     </p>
//     <p>
//         Best regards,<br />
//         MediaWatch team
//     </p>
// </div>
// `
// 	msgInvitationExistingAccount = `
// <div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
//     <p>
//         Hi there
//     </p>
//     <p>
//         {{.name}} has invited you to join the following team at <a href="https://mediawatch.io" target="_blank">MediaWatch</a>: {{.team}}.<br/>
//         Your 4-Digit Verification Code is: {{.nonce}}.
//     </p>
//     <p>
//         To accept or decline the invitation visit the following link: <a href="https://app.mediawatch.io/auth/invitation/{{.orgId}}/{{.memberId}}" target="_blank">Invitation</a>.<br/>
//     </p>
//     <p>
//         If you have any security or privacy discrimination concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
//     </p>
//     <p>
//         Best regards,<br />
//         MediaWatch team
//     </p>
// </div>
// `
// 	msgNewPass = `
// <div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
//     <p>
//         Hi %s
//     </p>
//     <p>
//         Your new password is: %s.
//     </p>
//     <p>
//         To login to your account follow this URL:<br/>
//         <a href="https://app.mediawatch.io/auth/login" target="_blank">https://app.mediawatch.io/auth/login</a>
//     </p>

//     <p>
//        Please make sure to change your password from your profile page once logged in.
//     </p>
//     <p>
//         Best regards,<br />
//         MediaWatch team
//     </p>
// </div>
// `

// 	msgReset = `
// <div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
//     <p>
//         Hi %s
//     </p>
//     <p>
//         You have requested a password reset for your MediaWatch account.<br/>
//         Your 4-Digit Verification Code is: %s.
//     </p>
//     <p>
//         To complete password reset follow this URL:<br/>
//         <a href="https://app.mediawatch.io/auth/reset/verify/%s" target="_blank">https://app.mediawatch.io/auth/reset/verify/%s</a>
//     </p>
//     <p>
//         The password reset link will expire in 24 hours.
//     </p>
//     <p>
//         If you didn't request a password reset or made it accidentially just ignore this message.<br/>
//         However, if you have any security concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
//     </p>
//     <p>
//         Best regards,<br />
//         MediaWatch team
//     </p>
// </div>
// `

// 	msgPin = `
// <div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
//     <p>
//         Hi %s
//     </p>
//     <p>
//         Your 4-Digit Verification Code is: %s.
//     </p>
//     <p>
//         To complete your registration follow this URL:<br/>
//         <a href="https://app.mediawatch.io/auth/verify/%s" target="_blank">https://app.mediawatch.io/auth/verify/%s</a>
//     </p>
//     <p>
//         Please confirm your account within the next 24 hours.
//     </p>
//     <p>
//         If you didn't register an account or made it accidentially just ignore this message.<br/>
//         However, if you have any security concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
//     </p>
//     <p>
//         Best regards,<br />
//         MediaWatch team
//     </p>
// </div>
// `

// 	msgDefault = `
// %s<br/><br/>

// Best regards,<br/>
// MediaWatch team<br/>
// `

// 	msgAccountDeletion = `
// <div style="font-family:'Segoe UI',Arial,Sans-Serif;font-size:10pt;">
//     <p>
//         Hi %s
//     </p>
//     <p>
//         We deactivated your account and we will complete the removal proccess within the next 14 days.
//     </p>
//     <p>
//         If you deleted your account by accident or if you have any security concerns please contact us immediately by email: <a href="mailto:press@mediawatch.io">press@mediawatch.io</a>.
//     </p>
//     <p>
//         Best regards,<br />
//         MediaWatch team
//     </p>
// </div>
// `
// )
