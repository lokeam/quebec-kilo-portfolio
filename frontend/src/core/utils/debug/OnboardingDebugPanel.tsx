import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Switch } from '@/shared/components/ui/switch';
import { Label } from '@/shared/components/ui/label';
import { getOnboardingDebugState, logDebugInfo } from './onboardingDebug';

/**
 * Onboarding Debug Panel
 *
 * A development-only component that provides easy toggles for testing
 * onboarding flow states. This should be removed in production.
 */
export function OnboardingDebugPanel() {
  const [isVisible, setIsVisible] = useState(false);
  const [debugState, setDebugState] = useState(getOnboardingDebugState());

  const updateDebugState = (key: keyof typeof debugState, value: boolean) => {
    const newState = { ...debugState, [key]: value };
    window.onboardingDebug = newState;
    setDebugState(newState);
    logDebugInfo('Debug Panel', `Updated ${key} to ${value}`);
  };

  const resetDebugState = () => {
    const defaultState = {
      bypassOnboarding: false,
      forceNewUser: false,
      forceIncompleteOnboarding: false,
      forceCompletedOnboarding: false,
      simulateProfileError: false,
      showDebugInfo: true,
    };
    window.onboardingDebug = defaultState;
    setDebugState(defaultState);
    logDebugInfo('Debug Panel', 'Reset all debug states');
  };

  if (!isVisible) {
    return (
      <div className="fixed bottom-4 right-4 z-50">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setIsVisible(true)}
          className="bg-yellow-100 border-yellow-300 text-yellow-800 hover:bg-yellow-200"
        >
          🔧 Debug
        </Button>
      </div>
    );
  }

  return (
    <div className="fixed bottom-4 right-4 z-50">
      <Card className="w-80 bg-yellow-50 border-yellow-300">
        <CardHeader className="pb-3">
          <CardTitle className="text-sm text-yellow-800 flex justify-between items-center">
            🔧 Onboarding Debug Panel
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsVisible(false)}
              className="h-6 w-6 p-0"
            >
              ×
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label htmlFor="bypass" className="text-xs">Bypass Onboarding</Label>
              <Switch
                id="bypass"
                checked={debugState.bypassOnboarding}
                onCheckedChange={(checked) => updateDebugState('bypassOnboarding', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="forceNew" className="text-xs">Force New User</Label>
              <Switch
                id="forceNew"
                checked={debugState.forceNewUser}
                onCheckedChange={(checked) => updateDebugState('forceNewUser', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="forceIncomplete" className="text-xs">Force Incomplete</Label>
              <Switch
                id="forceIncomplete"
                checked={debugState.forceIncompleteOnboarding}
                onCheckedChange={(checked) => updateDebugState('forceIncompleteOnboarding', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="forceComplete" className="text-xs">Force Complete</Label>
              <Switch
                id="forceComplete"
                checked={debugState.forceCompletedOnboarding}
                onCheckedChange={(checked) => updateDebugState('forceCompletedOnboarding', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="simulateError" className="text-xs">Simulate Error</Label>
              <Switch
                id="simulateError"
                checked={debugState.simulateProfileError}
                onCheckedChange={(checked) => updateDebugState('simulateProfileError', checked)}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label htmlFor="showDebug" className="text-xs">Show Debug Info</Label>
              <Switch
                id="showDebug"
                checked={debugState.showDebugInfo}
                onCheckedChange={(checked) => updateDebugState('showDebugInfo', checked)}
              />
            </div>
          </div>

          <div className="pt-2 border-t border-yellow-200">
            <Button
              variant="outline"
              size="sm"
              onClick={resetDebugState}
              className="w-full text-xs"
            >
              Reset All
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}