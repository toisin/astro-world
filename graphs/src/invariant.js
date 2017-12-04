// @flow
export default function invariant(v: any) {
  if (!v) {
    throw new Error('Invariant failed');
  }
}
