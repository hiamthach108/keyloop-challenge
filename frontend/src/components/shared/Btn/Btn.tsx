'use client';

import { Button, type ButtonProps } from 'antd';

export default function Btn({ children, ...props }: ButtonProps) {
  return <Button {...props}>{children}</Button>;
}
