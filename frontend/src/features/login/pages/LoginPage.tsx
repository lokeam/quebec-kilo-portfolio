import { useState } from 'react';
import { useAuthContext } from '@/core/auth/context-provider/AuthContext';

import QkoLogo from '@/shared/components/ui/LogoMap/brand/qko_logo';

export default function LoginPage() {
  const [authError, setAuthError] = useState<string | null>(null);
  const { login, isLoading } = useAuthContext();



  const handleLogin = async () => {
    try {
      await login();
    } catch (error) {
      setAuthError(error instanceof Error ? error.message : 'An error occurred');
      console.error('Login error:', error);
    }
  };

  console.error('Login Page authError:', authError);

  return (

    <div className="bk_login grid lg:grid-cols-2 h-screen">
      {/* ------  Login  ------ */}
      <div className="bk_login__cta bg-background text-foreground flex flex-col items-center justify-center px-4 py-6 sm:px-0 lg:py-0">
        <div className="max-w-md xl:max-w-xl">
          <h2 className="text-xl font-bold mb-8 text-foreground">Log in to your Q-KO Account</h2>

          {authError && (
            <div className="bg-destructive/10 border border-destructive text-destructive-foreground px-4 py-3 rounded relative mb-4" role="alert">
              <span className="block sm:inline">{authError}</span>
            </div>
          )}

          <form className="space-y-4 md:space-y-6 w-full max-w-md xl:max-w-xl">
            <button
              className="bg-primary text-primary-foreground hover:bg-primary/90 focus:ring-4 focus:outline-none focus:ring-ring font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center justify-center w-full border border-border hover:border-primary transition duration-500 ease-in-out"
              onClick={handleLogin}
              name="Login to Qko"
              type="button"
            >
              <QkoLogo className="w-4 h-4 me-2" />
              {isLoading ? 'Logging in...' : 'Create a QKO account'}
              </button>
          </form>
        </div>
      </div>

      {/* ------  Copy  ------ */}
      <div className="bk_login__mktg flex items-center justify-center bg-card text-card-foreground text-left p-6 sm:px-0 lg:py-0">
        <div className="max-w-xl px-6 xl:px-32">
          <div className="flex flex-row">
            <QkoLogo className="h-10 w-10 me-2" />
            <a className="font-domus text-foreground flex items-center text-4xl font-bold leading-none mb-4 hover:text-primary" href="/">Q-KO</a>
          </div>
          <h1 className="font-domus text-4xl font-extrabold mb-8 text-foreground">All of your games, in one place</h1>
          <p className="font-light leading-6 text-muted-foreground mb-4 lg:mb-8">By continuing, you acknowledge that you understand and agree to the <a className="text-primary hover:text-primary/80" href="/terms-and-privacy">Terms & Conditions</a> and <a className="text-primary hover:text-primary/80" href="/terms-and-privacy">Privacy Policy</a></p>
        </div>
      </div>
    </div>
  )
}
