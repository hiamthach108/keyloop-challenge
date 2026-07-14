module.exports = {
  '**/*.{ts,tsx}': () => 'npm run typecheck',
  '**/*.{ts,tsx,js,jsx,mjs,cjs}': (filenames) => [
    `npx eslint --fix ${filenames.join(' ')}`,
    `npx prettier --write ${filenames.join(' ')}`,
  ],
  '**/*.{md,json,css}': (filenames) => `npx prettier --write ${filenames.join(' ')}`,
};
