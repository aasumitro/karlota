import { Card } from '@/components/ui/card'
import AuthLayout from '../auth-layout'
import { UserLoginForm } from './components/user-login-form'
import {Link} from "@tanstack/react-router";

export default function SignIn() {
  return (
    <AuthLayout>
      <Card className='p-6'>
        <div className='flex flex-col space-y-2 text-left mb-4'>
          <h1 className='text-2xl font-semibold tracking-tight'>Login</h1>
          <p className='text-sm text-muted-foreground'>
            Enter your email and password below <br />
            to log into your account
          </p>
        </div>
        <UserLoginForm />
        <p className='mt-4 px-8 text-center text-sm text-muted-foreground'>
          Don't have an account?{' '}
          <Link
            to='/sign-up'
            className='underline underline-offset-4 hover:text-primary'
          >
            Sign up
          </Link>
          .
        </p>
      </Card>
    </AuthLayout>
  )
}