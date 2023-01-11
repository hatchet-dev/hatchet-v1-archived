import { atom } from "jotai";

const currOrgIdKey = "currOrgId";

const currOrgAtomInit = atom(localStorage.getItem(currOrgIdKey) ?? "");

export const currOrgAtom = atom(
  (get) => get(currOrgAtomInit),
  (get, set, newCurrOrgId: string) => {
    set(currOrgAtomInit, newCurrOrgId);
    localStorage.setItem(currOrgIdKey, newCurrOrgId);
  }
);
