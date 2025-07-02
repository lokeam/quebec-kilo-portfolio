import { useState } from 'react';
import { useAuth } from '@/core/auth/hooks/useAuth';

//import GoogleLogo from '@/shared/components/ui/LogoMap/brand/google_logo';
import QkoLogo from '@/shared/components/ui/LogoMap/brand/qko_logo';

export default function LoginPage() {
  const [authError, setAuthError] = useState<string | null>(null);
  const { login, isLoading } = useAuth();



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
      <div className="bk_login__cta bg-black flex flex-col items-center justify-center px-4 py-6 sm:px-0 lg:py-0">
        <div className="max-w-md xl:max-w-xl">
          <h2 className="text-xl font-bold mb-8">Log in to your Q-KO Account</h2>

          {authError && (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
              <span className="block sm:inline">{authError}</span>
            </div>
          )}

          <form className="space-y-4 md:space-y-6 w-full max-w-md xl:max-w-xl">
            <button
              className="text-white justify-center w-full border border-gray-700 hover:bg-gray-700 focus:ring-4 focus:outline-none focus:ring-[#4285F4]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#4285F4]/55 me-2 mb-2  transition duration-500 ease-in-out hover:border-vivid-blue dark:hover:border-vivid-blue"
              onClick={handleLogin}
              name="Login to Qko"
              type="button"
            >
              <QkoLogo className="w-4 h-4 me-2" />
              {isLoading ? 'Logging in...' : 'Continue with Auth0'}
              </button>
          </form>
        </div>
      </div>

      {/* ------  Copy  ------ */}
      <div className="bk_login__mktg flex items-center justify-center bg-white-smoke text-left p-6 sm:px-0 lg:py-0">
        <div className="max-w-xl px-6 xl:px-32">
          <div className="flex flex-row">
            <QkoLogo className="h-10 w-10 me-2" />
            <a className="font-domus text-charcoal flex items-center text-4xl font-bold leading-none mb-4 hover:text-strong-violet" href="/">Q-KO</a>
          </div>
          <h1 className="font-domus text-4xl font-extrabold mb-8 text-charcoal">All of your games, in one place</h1>
          <p className="font-light leading-6 text-charcoal mb-4 lg:mb-8">By continuing, you acknowledge that you understand and agree to the <a className="" href="/terms-and-privacy">Terms & Conditions</a> and <a className="" href="/terms-and-privacy">Privacy Policy</a></p>
        </div>
      </div>
    </div>
  )
}
