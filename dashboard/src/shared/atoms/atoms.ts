import { atom } from "jotai";
import { Organization, Team } from "shared/api/generated/data-contracts";

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

const currTeamKey = "currTeam";

const currTeamAtomInit = atom(getInitialValue(currTeamKey));

export const currTeamAtom = atom(
  (get) => get(currTeamAtomInit),
  (_get, set, newVal: Team) => {
    set(currTeamAtomInit, newVal);
    localStorage.setItem(currTeamKey, JSON.stringify(newVal));
  }
);
