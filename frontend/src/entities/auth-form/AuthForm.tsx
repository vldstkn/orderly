'use client';
import styles from './AuthForm.module.css';
import { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';

interface Inputs {
  name?: string;
  email: string;
  password: string;
}

export default function AuthForm() {
  const [isLogin, setIsLogin] = useState(true);
  const { register, handleSubmit, reset } = useForm<Inputs>();
  const onSubmit: SubmitHandler<Inputs> = data => {
    console.log(data);
    reset();
  };


  return <div className={styles.form_wrap}>
    <h2 className={styles.title}>Sign in or register</h2>
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
      <input placeholder={'Email'} type={'email'} className={styles.input} {...register('email')} />
      <input placeholder={'Password'} type={'password'} className={styles.input} {...register('password')} />
      {isLogin && <input placeholder={'Name'} className={styles.input} {...register('name')} />}
      <button className={styles.btn}>{isLogin ? 'Create' : 'Sign in'}</button>
    </form>
    <button className={styles.btn_toggle_type} onClick={() => setIsLogin(!isLogin)}>
      {isLogin ? 'Sign in?' : 'Create account?'}
    </button>
  </div>;
}