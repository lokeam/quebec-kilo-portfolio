import {
  Body,
  Column,
  Container,
  Head,
  Heading,
  Html,
  Img,
  Link,
  Preview,
  Row,
  Section,
  Hr,
  Text,
  Button,
} from '@react-email/components';

interface QKOBetaSignupProps {
  name: string;
  email: string;
  perks?: { id: number; description: string }[];
}

const defaultPerks = [
  {
    id: 1,
    description:
      'Exclusive Access: Gain early access to QKO and explore its features before the general public.',
  },
  {
    id: 2,
    description:
      'Special Perks: Free access to features in development such as wishlisted game deal pricing, Steam integration and other premium features available only to our beta testers.',
  },
  {
    id: 3,
    description:
      'Direct Influence: Your feedback will play a crucial role in refining and perfecting QKO!',
  },
];

export const QKOBetaSignup = ({
  name,
  email,
  perks = defaultPerks,
}: QKOBetaSignupProps) => (
  <Html>
    <Head />
    <Body style={main}>
      <Preview>Welcome to QKO Beta!</Preview>
      <Container style={container}>
        <Section style={logoContainer}>
          <Img
            src="/static/qko_logo.png"
            height="40"
            width="40"
            alt="QKO"
          />
        </Section>
        <Heading style={h1}>Welcome to QKO Beta!</Heading>
        <Text style={heroText}>
          Hi {name}, thanks for joining our early access list. We're excited to have you on board!
        </Text>

        <Hr style={hr} />

        <Text style={paragraph}>
          QKO is designed to be the only video game management platform that helps you keep track of the titles you own and the amount you spend over time.
        </Text>

        <Text style={paragraphEnd}>
          Here's what you can expect as a limited beta participant:
        </Text>

        {/* Email-compatible list using Text components */}
        {perks.map((perk, index) => (
          <Text key={perk.id} style={paragraph}>
            • {perk.description}
          </Text>
        ))}

        <Section style={buttonContainer}>
          <Button style={button} href="https://dashboard.q-ko.com/login">
            Sign in
          </Button>
        </Section>

        <Hr style={hr} />

        <Text style={paragraph}>
          We are thrilled to have you on board and look forward to your insights! If you have any questions or need further information, please do not hesitate to reach out.
        </Text>

        <Text style={paragraph}>— The QKO team</Text>

        <Section>
          <Text style={footerText}>
            ©2025 QKO. All rights reserved. <br />
            <br />
            <Link href="https://q-ko.com" style={footerLink}>
              Visit QKO
            </Link>
          </Text>
        </Section>
      </Container>
    </Body>
  </Html>
);

QKOBetaSignup.PreviewProps = {
  name: 'John Doe',
  email: 'john.doe@example.com',
  perks: defaultPerks,
} as QKOBetaSignupProps;

export default QKOBetaSignup;

const button = {
  backgroundColor: '#3D6FB7',
  borderRadius: '5px',
  color: '#fff',
  fontSize: '16px',
  fontWeight: 'bold',
  textDecoration: 'none',
  textAlign: 'center' as const,
  display: 'block',
  width: '200px',
  padding: '10px',
  margin: '0 auto',
};

const buttonContainer = {
  textAlign: 'center' as const,
  marginTop: '20px',
};

const listItem = {
  color: '#000',
  fontSize: '14px',
  lineHeight: '24px',
  marginBottom: '12px',
  paddingLeft: '0',
};

const footerText = {
  fontSize: '12px',
  color: '#b7b7b7',
  lineHeight: '15px',
  textAlign: 'left' as const,
  marginBottom: '50px',
};

const footerLink = {
  color: '#b7b7b7',
  textDecoration: 'underline',
};

const main = {
  backgroundColor: '#ffffff',
  margin: '0 auto',
  fontFamily:
    "-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif",
};

const container = {
  margin: '0 auto',
  padding: '0px 20px',
};

const logoContainer = {
  marginTop: '32px',
};

const h1 = {
  color: '#1d1c1d',
  fontSize: '36px',
  fontWeight: '700',
  margin: '30px 0',
  padding: '0',
  lineHeight: '42px',
};

const heroText = {
  fontSize: '20px',
  lineHeight: '28px',
  marginBottom: '10px',
};

const hr = {
  borderColor: '#e6ebf1',
  margin: '20px 0',
};

const text = {
  color: '#000',
  fontSize: '14px',
  lineHeight: '24px',
};

const paragraph = {
  color: '#000',
  fontSize: '16px',
  lineHeight: '24px',
  textAlign: 'left' as const,
};

const paragraphEnd = {
  color: '#000',
  fontSize: '16px',
  lineHeight: '24px',
  textAlign: 'left' as const,
  marginBottom: '20px',
};