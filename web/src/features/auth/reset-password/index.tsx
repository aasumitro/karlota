import { Card } from '@/components/ui/card'
import AuthLayout from '../auth-layout'
import { ResetPasswordForm } from './components/reset-password-form'

export default function ResetPassword() {
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