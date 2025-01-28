import { useEffect, useState } from 'react';

// ShadCN UI Components
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/components/ui/card';
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from '@/shared/components/ui/dialog';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';

// Icons
import { Trash2 } from "lucide-react";
import { TriangleAlert } from "lucide-react";


export function DangerZoneSection() {
  const [confirmText, setConfirmText] = useState<string>('');
  const [isDeleteEnabled, setIsDeleteEnabled] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setIsDeleteEnabled(confirmText.toLowerCase() === "delete my account")
  }, [confirmText]);

  const handleDelete = () => {
    console.log("Account delete action triggered.")
    // setIsLoading(true)
    // TODO: Implement actual account deletion logic
  }

  return (
    <Card className="border-red-500">
      <CardHeader>
        <CardTitle>Danger Zone</CardTitle>
        <CardDescription>Delete your account.</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-4">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">
              Once you delete your account, there is no going back.
            </p>
            <p className="text-sm text-gray-400">Please be certain.</p>
          </div>

          <Dialog>
            <DialogTrigger asChild>
              <Button
                className="bg-gray-200 h-11 border-red-500 text-red-500 hover:text-white hover:bg-red-600 hover:border-red-600 focus:ring-red-900 dark:bg-gray-800 dark:hover:bg-red-600 dark:border-2 transition duration-500 ease-in-out"
                variant="outline"
              >
                Delete Account
              </Button>
            </DialogTrigger>

            <DialogContent className="sm:max-w-[425px] space-y-0 gap-0">
              <DialogHeader>
                <DialogTitle className="text-2xl text-center font-semibold text-red-600 pb-2">
                  Warning
                </DialogTitle>
                <DialogDescription hidden>
                  Make changes to your profile here. Click save when you're done.
                </DialogDescription>
              </DialogHeader>

              <div className="flex flex-col items-center justify-center">
                <div className="flex flex-row items-center justify-center">
                  <TriangleAlert className="w-10 h-10 text-yellow-500" />
                </div>
                <div className="flex flex-col items-center justify-center space-y-4">
                  <h3 className="flex text-center items-center justify-center text-md mt-3">
                    Are you sure that you want to delete your account?
                  </h3>
                  <p className="font-bold flex items-center justify-center">
                    This action cannot be undone.
                  </p>
                  <p className="text-s text-gray-600 dark:text-gray-400">
                    To confirm, type{" "}
                    <span className="font-bold text-red-800 dark:text-red-500">
                      delete my account
                    </span>{" "}
                    below.
                  </p>

                  <Input
                    type="text"
                    value={confirmText}
                    onChange={(e) => setConfirmText(e.target.value)}
                    className="border-gray-700/60 bg-transparent border-2 rounded-md p-2 mb-2"
                    placeholder="Type 'delete my account'"
                  />
                </div>
              </div>

              <DialogFooter className="flex flex-col items-center justify-center">
                <div className="flex flex-col items-center justify-center w-full space-y-3 mt-10">
                  <DialogClose asChild>
                    <Button
                      className="flex flex-row justify-between items-center border-gray-700/60 mr-1 w-full mb-3 lg:mb-0 transition duration-500 ease-in-out hover:border-blue-500 hover:bg-blue-500 hover:text-white"
                    >
                      I changed my mind, get me out of here!
                    </Button>
                  </DialogClose>
                  <Button
                    type="button"
                    onClick={handleDelete}
                    disabled={!isDeleteEnabled}
                    className={`bg-transparent flex flex-row justify-between items-center mr-1 w-full border-red-500 text-red-500 hover:text-white hover:bg-red-800 focus:ring-red-800 hover:border-red-800 dark:hover:bg-red-800 transition duration-500 ease-in-out disabled:pointer-events-none disabled:border-gray-700/60 disabled:text-gray-600 dark:disabled:text-gray-400`}
                    variant="outline"
                  >
                    <span>{isLoading ? "Deleting..." : "Yes, I want to delete my account"}</span>
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </div>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </CardContent>
    </Card>
  )
}
