import initialData from "./data.json" with { type: "json" };

type ChairJsonType = { id: string; name: string; model: string; token: string };
type OwnerJsonType = {
  id: string;
  name: string;
  token: string;
  chairs: ChairJsonType[];
};
type JsonType = { owners: OwnerJsonType[] };

const jsonData = initialData satisfies JsonType;

export const getOWners = () => {
  return jsonData.owners.map((owner) => ({
    id: owner.id,
    name: owner.name,
    token: owner.token,
  }));
};

export const getChairs = (ownerId: OwnerJsonType["id"]) => {
  return jsonData.owners.find((owner) => owner.id === ownerId)?.chairs ?? [];
};
