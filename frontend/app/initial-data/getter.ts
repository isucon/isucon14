import initialData from "./data.json" with { type: "json" };

type ChairJsonType = { id: string; name: string; model: string; token: string };
type OwnerJsonType = { id:string; name: string; token: string; chairs: ChairJsonType[] };
type JsonType = { owners: OwnerJsonType[] };

const isJsonType = (value: any): value is JsonType => {
  return value?.owners?.every(
    (owner: Partial<OwnerJsonType> | undefined) =>
      typeof owner?.name === "string" &&
      typeof owner?.token === "string" &&
      owner?.chairs?.every(
        (chair: Partial<ChairJsonType> | undefined) =>
          typeof chair?.id === "string" &&
          typeof chair?.name === "string" &&
          typeof chair?.model === "string" &&
          typeof chair?.token === "string",
      ),
  );
};

const jsonData = (() => {
  const json = initialData;
  if (isJsonType(json)) {
    return json;
  } else {
    throw Error("type of initalData isn't `JsonType` ");
  }
})();


export const getOWners = () => {
  return jsonData.owners.map((owner) => ({id: owner.id, name: owner.name, token: owner.token}))
}

export const getChairs = (ownerId: OwnerJsonType["id"]) => {
  return jsonData.owners.find((owner) => owner.id === ownerId)?.chairs ?? []
}