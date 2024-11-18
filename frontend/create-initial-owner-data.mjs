/**
 * 初期データが作成されるまで一旦こちらで作成してdata.jsonに突っ込むようにする
 */

import { writeFileSync } from "fs";

const  getCookieValue = (cookieString, cookieName) => {
  const regex = new RegExp(`(?:^|; )${cookieName}=([^;]*)`);
  const match = cookieString.match(regex);
  return match ? match[1] : null; // 値が見つかった場合は返す、見つからなければ null
}

const create = async() => {
  const candidates = [
    {name: "オーナー1", chairs: [{
      name: "chair1", model: "model1"
    }, {
      name: "chair2", model: "model2"
    }, {
      name: "chair3", model: "model3"
    }]},
    {name: "オーナー2", chairs: [
      {
        name: "chair1", model: "model1"
      }
    ]},
    {name: "オーナー3", chairs: []}
  ]

  const BASE_URL = "http://localhost:8080";

  const candidate = await Promise.all(candidates.map(async (candidate) => {
    const ownerFetch = await fetch(`${BASE_URL}/api/owner/owners`, {
      "body": JSON.stringify({
        name: candidate.name
      }),
      "credentials": 'include',
      "method": "POST"
    })
    /**
     * @type {{id: string,chair_register_token: string}}
     */
    const json = await ownerFetch.json();

    const chairRegisterToken = json.chair_register_token;

    // set-cookie ヘッダーを取得
    const cookie = ownerFetch.headers.get('set-cookie');

    // owner_session クッキーを探す
    const ownerSessionValue = getCookieValue(cookie, "owner_session");
    
    const chairs = await Promise.all(candidate.chairs.map( async (chair) => {
      const chairFetch = await fetch(`${BASE_URL}/api/chair/chairs`, {
        body: JSON.stringify({name: chair.name, model: chair.model, chair_register_token: chairRegisterToken}),
        method: "POST",
        "credentials": 'include',
      })
      let json;
      /**
       * @type {{id: string, owner_id: string}}
       */
      try {
       json = await chairFetch.json()
      }catch(e) {
        console.error(e)
      }

      const cookie =  chairFetch.headers.get('set-cookie');
      const chairSessionValue = getCookieValue(cookie, "chair_session");
      return {
        id: json.id,
        name: chair.name,
        model: chair.model,
        token: chairSessionValue
      }
    }))

    return {
      id: json.id,
      name: candidate.name,
      token: ownerSessionValue,
      chairs
    }
  }))
  writeFileSync("./initial-owner-data.json", JSON.stringify({owners: candidate}, null, 2))
  console.log('created')
}

create()