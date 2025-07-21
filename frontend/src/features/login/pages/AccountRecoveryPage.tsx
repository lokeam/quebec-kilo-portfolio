import React, { useState } from 'react';
import { BrandLogo } from '@/features/navigation/organisms/SideNav/BrandLogo';

export default function AccountRecoveryPage() {
  const [email, setEmail] = useState('');
  const [acceptedTerms, setAcceptedTerms] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!acceptedTerms) {
      setError('You must accept the terms and conditions.');
      return;
    }
    setIsSubmitting(true);
    // Placeholder for submit logic
    setTimeout(() => {
      setIsSubmitting(false);
    }, 1000);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-black px-4">
      <div className="w-full max-w-md bg-[#232734] rounded-2xl shadow-lg p-8 md:p-10 flex flex-col items-center">
        <div className="mb-6 flex flex-col items-center">
          <BrandLogo />
        </div>
        <h2 className="text-2xl font-bold text-white mb-2 text-center">Account Recovery</h2>
        <p className="text-gray-400 mb-6 text-center">Your account is marked for deletion.</p>
        <p className="text-gray-400 mb-6 text-center">If you would like to reactivate it, please enter the email address associated with your account and we will complete the process.</p>
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4 w-full text-center" role="alert">
            <span className="block sm:inline">{error}</span>
          </div>
        )}
        <form className="w-full space-y-4" onSubmit={handleSubmit}>
          <div className="text-left">
            <label htmlFor="email" className="block text-white font-semibold mb-1">Your email</label>
            <input
              id="email"
              type="email"
              className="w-full rounded-lg border border-gray-700 bg-[#232734] text-gray-200 placeholder-gray-500 px-4 py-3 focus:outline-none focus:ring-2 focus:ring-vivid-blue"
              placeholder="name@example.com"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="flex items-center">
            <input
              id="terms"
              type="checkbox"
              checked={acceptedTerms}
              onChange={e => setAcceptedTerms(e.target.checked)}
              className="w-4 h-4 text-vivid-blue bg-gray-100 border-gray-300 rounded focus:ring-vivid-blue focus:ring-2"
            />
            <label htmlFor="terms" className="ml-2 text-gray-400 text-sm">
              I no longer want to delete my account.
            </label>
          </div>
          <button
            type="submit"
            className="w-full bg-vivid-blue hover:bg-blue-700 text-white font-medium rounded-lg text-sm px-5 py-3 transition duration-300 disabled:opacity-50"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Sending...' : 'Reactivate account'}
          </button>
        </form>
      </div>
    </div>
  );
}
