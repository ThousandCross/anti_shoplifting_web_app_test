import { initializeApp } from "https://www.gstatic.com/firebasejs/9.9.3/firebase-app.js";
//import { getMessaging } from "https://www.gstatic.com/firebasejs/9.9.3/firebase-messaging.js";
import { getMessaging, onBackgroundMessage, isSupported } from "https://www.gstatic.com/firebasejs/9.9.3/firebase-messaging-sw.js";
//var ASSETS = ['/static/css/toastifier.min.css', '/js/index.js', '/style/style.css'];
// var ASSETS = ['./static/css/toastifier.min.css'];
// import { toastifier } from "./static/js/toastifier.js";
//importScripts('./static/js/toastifier.js');
// import './static/css/toastifier.min.css';
// import toastifier from 'toastifier';



// self.oninstall = function (evt) {
//   evt.waitUntil(caches.open('offline-cache-name').then(function (cache) {
//     return Promise.all(ASSETS.map(function (url) {
//       return fetch(url).then(function (response) {
//         return cache.put(url, response);
//       });
//     }));
//   }))
// };

// Your web app's Firebase configuration
const firebaseConfig = {
    apiKey: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    authDomain: "anti-shoplifting-dev.firebaseapp.com",
    projectId: "anti-shoplifting-dev",
    storageBucket: "anti-shoplifting-dev.appspot.com",
    messagingSenderId: "916044113120",
    appId: "1:916044113120:web:3bce2a9f9dde897f829c85",
    measurementId: "G-HJR18QLVZF"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
isSupported()
    .then(() => {
        const messaging = getMessaging(app)
        onBackgroundMessage(messaging, (payload) => {
            console.log('[firebase-messaging-sw.js] Received background message ', payload);

            const notificationData = payload.notification;
            if (!notificationData) return;
            const { title, body } = notificationData;
            if (!(title && body)) return;

            const notificationOptions = {
                body: 'Background Message body.',
                //icon: '/firebase-logo.png'
                icon: undefined // サーバーのpublicフォルダーに入れている画像を通知に入れることができます。
            };

            // Customize notification here
            const notificationTitle = 'Background Message Title';
            self.registration.showNotification(notificationTitle,
                notificationOptions);
        })
    }).catch(
        console.log("")
        );

addEventListener('fetch', event => {
  //console.log('event', event);
  event.waitUntil(async function() {
    // Exit early if we don't have access to the client.
    // Eg, if it's cross-origin.
    if (!event.clientId) return;
        
    // Get the client.
    const client = await clients.get(event.clientId);
    // Exit early if we don't get the client.
    // Eg, if it closed.
    if (!client) return;
        
    // Send a message to the client.
    //console.log('Send a message to the client.');
    client.postMessage({
      msg: "Hey I just got a fetch from you!",
      url: event.request.url
    });
}());
});

//プッシュ通知が行われると「push」イベントが起動する
self.addEventListener("install", function(event) {
    self.skipWaiting();
    console.log("Installed", event);

    // event.waitUntil(caches.open('offline-cache-name').then(function (cache) {
    //   return Promise.all(ASSETS.map(function (url) {
    //     return fetch(url).then(function (response) {
    //       return cache.put(url, response);
    //     });
    //   }));
    // }))
});

self.addEventListener("activate", function(event) {
    console.log("Activated", event);
});

self.addEventListener('push', function (event) {
    //console.log('Received a push message', event);
    var title = "万引き発生！";
    var body = "万引きが発生しました！！";

    event.waitUntil(
        self.registration.showNotification(title, {
            body: body,
            icon: '/static/img/Anti_Shoplifting_Dev_Icon.png',
            tag: 'push-notification-tag'
        })
    );
});

self.addEventListener("message", function(event) {
  //event.source.postMessage("Responding to " + event.data);
  self.clients.matchAll().then(all => all.forEach(client => {
      client.postMessage("Responding to " + event.data);
  }));
});

self.addEventListener("pushsubscriptionchange", (event) => {
    const subscription = swRegistration.pushManager
      .subscribe(event.oldSubscription.options)
      .then((subscription) =>
        fetch("register", {
          method: "post",
          headers: {
            "Content-type": "application/json",
          },
          body: JSON.stringify({
            endpoint: subscription.endpoint,
          }),
        }),
      );
    event.waitUntil(subscription);
  }, false);

self.addEventListener('notificationclick', function (event) {
    event.notification.close();
    clients.openWindow("/incidents");
}, false);