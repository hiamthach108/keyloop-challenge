export default function Header() {
  return (
    <header className="border-b border-neutral-200 bg-white">
      <div className="mx-auto flex min-h-16 w-full max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
        <div>
          <p className="text-base font-semibold text-neutral-950">Intelligent Inventory</p>
          <p className="text-xs text-neutral-500">Dealership stock overview</p>
        </div>
        <span className="rounded-md bg-blue-50 px-3 py-1 text-xs font-medium text-blue-700">
          Scenario B
        </span>
      </div>
    </header>
  );
}
