import { useState } from 'react';
import { useAuthContext } from '@/core/auth/context-provider/AuthContext';
import QkoLogo from '@/shared/components/ui/LogoMap/brand/qko_logo';

export default function DeletedAccountPage() {
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

  return (
    <div className="min-h-screen flex items-center justify-center bg-black px-4">
      <div className="w-full max-w-md bg-[#232734] rounded-2xl shadow-lg p-8 md:p-10 flex flex-col items-center">
        <div className="mb-6 flex flex-col items-center">
          <QkoLogo className="w-4 h-4 me-2" />
        </div>
        <h2 className="text-2xl font-bold text-white mb-2 text-center">Account Deleted</h2>
        {authError && (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
              <span className="block sm:inline">{authError}</span>
            </div>
          )}

        <p className="text-gray-400 mb-6 text-center">Your account was marked for deletion and has been permanently deleted.</p>
        <p className="text-gray-400 mb-6 text-center">Please create a new account if you would like to continue using QKO.</p>

        <form className="w-full space-y-4">
        <button
              className="text-white justify-center w-full border border-gray-700 hover:bg-gray-700 focus:ring-4 focus:outline-none focus:ring-[#4285F4]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#4285F4]/55 me-2 mb-2  transition duration-500 ease-in-out hover:border-vivid-blue dark:hover:border-vivid-blue"
              onClick={handleLogin}
              name="Login to Qko"
              type="button"
            >
              {isLoading ? 'Logging in...' : 'Create a QKO account'}
              </button>
        </form>
      </div>
    </div>
  );
}
