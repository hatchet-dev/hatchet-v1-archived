import { atom } from "jotai";
import { Organization } from "shared/api/generated/data-contracts";

const getInitialValue = (key: string) => {
  const item = localStorage.getItem(key);

  if (item !== null) {
    return JSON.parse(item);
  }

  return null;
};

const currOrgKey = "currOrg";

const currOrgAtomInit = atom(getInitialValue(currOrgKey));

export const currOrgAtom = atom(
  (get) => get(currOrgAtomInit),
  (_get, set, newVal: Organization) => {
    set(currOrgAtomInit, newVal);
    localStorage.setItem(currOrgKey, JSON.stringify(newVal));
  }
);
