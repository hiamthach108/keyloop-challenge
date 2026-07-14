import nextVitals from 'eslint-config-next/core-web-vitals';
import nextTypescript from 'eslint-config-next/typescript';
import prettier from 'eslint-config-prettier';

const config = [
  ...nextVitals,
  ...nextTypescript,
  prettier,
  {
    ignores: ['.next/**', 'node_modules/**', 'storybook-static/**'],
  },
];

export default config;
