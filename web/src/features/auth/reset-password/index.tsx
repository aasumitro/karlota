import { Card } from '@/components/ui/card'
import AuthLayout from '../auth-layout'
import { ResetPasswordForm } from './components/reset-password-form'

export default function ResetPassword() {
  const queryParams = new URLSearchParams(location.search);
  const exchangeToken = queryParams.get('exchange_token');

  if (!exchangeToken) {
    return (
      <div className='h-svh'>
        <div className='m-auto flex h-full w-full flex-col items-center justify-center gap-2'>
          <h1 className='text-[1rem] font-semibold tracking-wide'>
            Exchange Token is Required
          </h1>
          <p className="text-center">
            Please check again your email and validate
            <br/> the link that we've sent or request new {" "}
            <a href="/forgot-password" className="border-b-2 border-black">reset link</a>
          </p>
        </div>
      </div>
    );
  }

  return (
    <AuthLayout>
      <Card className='p-6'>
        <div className='mb-2 flex flex-col space-y-2 text-left'>
          <h1 className='text-md font-semibold tracking-tight'>
            Create new password
          </h1>
          <p className='text-sm text-muted-foreground'>
            Please enter your new password below. <br />
            Make sure it's something secure that you'll remember.
          </p>
        </div>
        <ResetPasswordForm />
      </Card>
    </AuthLayout>
  )
}