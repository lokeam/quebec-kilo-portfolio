import type { APIRoute } from 'astro';
import { disposibleDomains } from '@/pages/api/disposible-domains';
import { Resend } from 'resend';
import { render } from '@react-email/render';
import QKOBetaSignup from '@/pages/api/templates/qko-beta-signup';

const isValidEmail = (email: string): boolean => {
  return /^[^@\\s]+@[^@\\s]+\\.[^@\\s]+$/.test(email) && email.length <= 254;
};

const isDisposableEmail = (email: string): boolean => {
  const domain = email.split('@')[1]?.toLowerCase();
  if (!domain) return false;

  const subdomain = domain.split('.')[0];
  return disposibleDomains.has(subdomain);
};

const isValidSubmission = (email: string, name: string): boolean => {
  return isValidEmail(email) && !isDisposableEmail(email) &&
         name.length >= 2 && name.length <= 100;
};


export const POST: APIRoute = async ({ request }) => {

  const form = await request.formData();
  const fullName = form.get('fullName')?.toString() || '';
  const email = form.get('email')?.toString().trim().toLowerCase() || '';
  const honeypot = form.get('website')?.toString() || '';

  // Honeypot check
  if (honeypot) {
    return new Response(JSON.stringify({ ok: true }), { status: 200 });
  }

  // Basic validation
  if (!isValidSubmission(email, fullName)) {
    console.log('Invalid submission:', { email, fullName });
    return new Response(JSON.stringify({ ok: true }), { status: 200 });
  }

  // Send email notification
  try {
    const resend = new Resend(import.meta.env.RESEND_API_KEY);

    // Render the React Email Beta Signup template
    const betaSignupEmailTemplate = await render(
      QKOBetaSignup({
        name: fullName,
        email: email,
      })
    );

    // Send email via Resend
    await resend.emails.send({
      from: import.meta.env.EMAIL_FROM,
      to: [email], // Send to the user who signed up
      subject: 'Welcome to QKO Beta!',
      html: betaSignupEmailTemplate,
    });

    // Also send notification to admin
    const adminNotification = `
      <h2>New Beta Access Signup</h2>
      <p><strong>Name:</strong> ${fullName}</p>
      <p><strong>Email:</strong> ${email}</p>
      <p><strong>Time:</strong> ${new Date().toISOString()}</p>
    `;

    await resend.emails.send({
      from: import.meta.env.EMAIL_FROM,
      to: [import.meta.env.BETA_SIGNUP_EMAIL],
      subject: 'New Beta Signup',
      html: adminNotification,
    });

    console.log('Email sent successfully:', { email, fullName });
    return new Response(JSON.stringify({ ok: true }), { status: 200 });
  } catch (error) {
    console.error('Email sending error:', error);
    return new Response(JSON.stringify({ ok: true }), { status: 200 });
  }

}