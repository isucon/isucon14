type ChairJsonType = { id: string; name: string; model: string; token: string };
type OwnerJsonType = {
  id: string;
  name: string;
  token: string;
  chairs: ChairJsonType[];
};
type JsonType = { owners: OwnerJsonType[] };


const initialOwnerData = __INITIAL_OWNER_DATA__ || { owners: []} satisfies JsonType;

export const getOWners = () => {
  return initialOwnerData.owners.map((owner) => ({
    id: owner.id,
    name: owner.name,
    token: owner.token,
  }));
};

export const getChairs = (ownerId: OwnerJsonType["id"]) => {
  return initialOwnerData.owners.find((owner) => owner.id === ownerId)?.chairs ?? [];
};
