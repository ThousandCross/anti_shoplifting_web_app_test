//モーダル終了時のイベント登録
document.getElementById('staticBackdrop').addEventListener('hidden.bs.modal', function (event) {
    //tbodyの子要素を全て削除
    while(document.getElementById('tbody').firstChild){
        document.getElementById('tbody').removeChild(tbody.firstChild);
    }
})

let search = document.getElementById('search');
search.addEventListener('click', () => {
    let api = 'https://zipcloud.ibsnet.co.jp/api/search?zipcode=';
    let zipcode = document.getElementById('zipcode');
    let prefecture = document.getElementById('prefecture');
    let city = document.getElementById('city');
    let param = zipcode.value.replace("-", ""); //入力された郵便番号から「-」を削除
    let url = api + param;
    console.log(url);

    fetchJsonp(url, {
      timeout: 10000, //タイムアウト時間
    })
      .then((response) => {
        //error.textContent = ''; //HTML側のエラーメッセージ初期化
        return response.json();
      })
      .then((data) => {
        if (data.status === 200) {
          if (data.results.length > 1) {
            // 住所を2件以上抽出した場合の処理
            // モーダル取得
            var myModal = new bootstrap.Modal(document.getElementById("staticBackdrop"), {});

            // 作成したHTML要素をtbody要素に追加する
            var tbody = document.getElementById('tbody');
            for(var i=0; i<data.results.length; i++) {
              // tr
              var tr_element = document.createElement('tr');
              // th
	            var th_element = document.createElement('th');
              th_element.scope = "row";
	            th_element.textContent = i+1;
              // td
              var td_element = document.createElement('td');
              td_element.textContent = data.results[i].address1 + data.results[i].address2 + data.results[i].address3;
              
              function setAddress(e){
                let prefecture_id = Number(this.data.results[this.index].prefcode)
                // id="city"の入力フォームを上書き
                document.getElementById('prefecture').options[prefecture_id].selected = true;
                document.getElementById('city').value = this.data.results[this.index].address2 + this.data.results[this.index].address3;
                //var myModal = new bootstrap.Modal(document.getElementById("staticBackdrop"), {});
                //tbodyの子要素を全て削除
                while(document.getElementById('tbody').firstChild){
                  document.getElementById('tbody').removeChild(tbody.firstChild);
                }
                this.modal.hide();
                //myModal.close();
              };

              td_element.addEventListener('click', {modal: myModal, data: data, index: i, handleEvent: setAddress});
              // td要素にクリックイベント登録

	            tr_element.appendChild(th_element);
              tr_element.appendChild(td_element);
              tbody.appendChild(tr_element);
            }
            myModal.show();
            //document.getElementById("staticBackdrop").style.visibility ="visible";
          } else {
            // 住所を1件のみ抽出した場合の処理
            let prefecture_id = Number(data.results[0].prefcode)
            // id="city"の入力フォームを上書き
            prefecture.options[prefecture_id].selected = true;
            city.value = data.results[0].address2 + data.results[0].address3;
          }

        } else if (data.status === 400) { //エラー時
          alert('該当する住所が見つかりませんでした');
        } else if (data.results === null) {
          alert('該当する住所が見つかりませんでした');
        } else {
          alert('該当する住所が見つかりませんでした');
        }
    })
    .catch((ex) => { //例外処理
        console.log(ex);
    });
}, false);