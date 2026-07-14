import Btn from '@/components/shared/Btn/Btn';

export default function HomeHero() {
  return (
    <section className="flex flex-col gap-4">
      <p className="text-sm font-medium uppercase tracking-wide text-neutral-500">
        Dreon Next.js Boilerplate
      </p>
      <div className="flex max-w-2xl flex-col gap-3">
        <h1 className="text-4xl font-semibold text-neutral-950">Build the product first.</h1>
        <p className="text-base leading-7 text-neutral-600">
          A modern App Router starter with Next.js, React, Ant Design, Tailwind CSS,
          ESLint, Prettier, and Storybook.
        </p>
      </div>
      <div>
        <Btn type="primary">Get started</Btn>
      </div>
    </section>
  );
}
