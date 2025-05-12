// Import the functions you need from the SDKs you need
import { initializeApp } from "https://www.gstatic.com/firebasejs/9.9.3/firebase-app.js";
import { getMessaging, getToken, onMessage } from "https://www.gstatic.com/firebasejs/9.9.3/firebase-messaging.js";
// import 'toastifier/dist/toastifier.min.css';
// import toastifier from 'toastifier';
//import { toastifier } from "./static/js/toastifier.js";
//importScripts('./static/js/toastifier.js');
//import toastifier from './toastifier.js';
// import '../../node_modules/toastifier/dist/toastifier.min.css';
// import toastifier from '../../node_modules/toastifier';
var head  = document.getElementsByTagName('head')[0];
var link  = document.createElement('link');
//link.id   = cssId;
link.rel  = 'stylesheet';
link.type = 'text/css';
link.href = '/static/css/toastifier.min.css';
link.media = 'all';
head.appendChild(link);
import {toastifier} from "./toastifier.js";

// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries
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
export const messaging = getMessaging(app);

//serviceWorker registeration
if ('serviceWorker' in navigator) {
    window.addEventListener('load', function () {
        navigator.serviceWorker.register('/firebase-messaging-sw.js', { type: 'module' })
            .then(function (registration) {
                // Service worker registration done
                //console.log('Registration Successful', registration);
                console.log('Registration Successful');

                Notification.requestPermission()
                    .then((permission) => {
                        if (permission == 'granted') {
                            // 許可
                            console.log("granted!");
                            // Get registration token. Initially this makes a network call, once retrieved
                            // subsequent calls to getToken will return from cache.
                            // const messaging = getMessaging(firebase_app);
                            getToken(messaging, {
                                serviceWorkerRegistration: registration,
                                vapidKey: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxc",
                            }).then((currentToken) => {
                                if (currentToken) {
                                    // Send the token to your server and update the UI if necessary
                                    // alert(currentToken);
                                    console.log("currentToken:", currentToken);

                                    // サブスクリプション登録を行う
                                    navigator.serviceWorker.ready.then(p => {
                                        console.log("serviceWorker.ready!");
                                        p.pushManager.getSubscription().then(function(subscription) {
                                            // すでに購読済みの通知があるかを判定。
                                            //if (subscription) return subscription;
                                            //console.log("get new subscription!");
                                            // if(!pushSubscription){
                                            // 通知の購読が存在しない場合は登録する。
                                            let re = p.pushManager.subscribe({
                                                userVisibleOnly: true,
                                                applicationServerKey: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
                                            });
                                            console.log("subscribed!");
                                            // } else {
                                            //     console.log("not subscribed!");
                                            // }
                                            if (location.pathname == '/') {
                                                var csrf_token = document.getElementById("csrf_token").value;
                                                console.log(csrf_token);
                                                // form を動的に生成
                                                var form = document.createElement('form');
                                                let formData = new FormData(form);
                                                // {{.CSRFToken}} not work, but  csrf_token worked!
                                                //formData.append("csrf_token", "{{.CSRFToken}}");
                                                formData.append("csrf_token", csrf_token);
                                                formData.append("fcm_token", currentToken);

                                                fetch('/register_fcm_token', {
                                                    method: "post",
                                                    body: formData,
                                                })
                                                .then(response => response.json())
                                                .then(data => {
                                                    console.log(data);
                                                    console.log(data.result);
                                                    console.log(data.message);
                                                })
                                            }
                                        });
                                    });
                                } else {
                                    // Show permission request UI
                                    //console.log('No registration token available. Request permission to generate one.');
                                    alert('No registration token available. Request permission to generate one.');
                                }
                            }).catch((err) => {
                                //console.log('An error occurred while retrieving token.', err);
                                alert('An error occurred while retrieving token.');
                            });
                        } else if (permission == 'denied') {
                            // 拒否
                            console.log("denied!");
                            return;
                        } else if (permission == 'default') {
                            // 無視
                            console.log("default.ignored!");
                            return;
                        }
                    });
            }).catch(function(error) {
                console.log('Registration Failed', error);
            });

            navigator.serviceWorker.addEventListener("message", (event) => {
                // console.log('navigator.serviceWorker.addEventListener called');
                // console.log("Got reply from service worker: ",  event);
                if (event.data.messageType) {
                    // pushイベントの時は、event.data.messageType == 'push-received'になる
                    if (event.data.messageType == 'push-received') {
                        // console.log("push event received in initialization");
                        // event.data.data以降はAPIにて定義済
                        if(event.data.data.type == "incidents") {
                            toastifier('万引きが発生しました!!',{
                                duration: 5000,
                                showIcon: true,
                                type: 'warn',
                                animation: 'bounce',
                              });
                        } else if (event.data.data.type == "sensors") {
                            var label_sensor_cd = document.getElementById("dashboard-weight-change-sensor-cd");
                            var label_type = document.getElementById("dashboard-weight-change-type");
                            var label_weight_value = document.getElementById("dashboard-weight-change-value");

                            label_sensor_cd.innerHTML =  "Sensor ID : " + String(parseInt(event.data.data.sensor_cd));
                            var weight_value_diff = event.data.data.weight_value_diff;
                            if ( parseFloat(weight_value_diff) > 0) {
                                label_type.style.borderTop = "0px solid transparent";
                                label_type.style.borderBottom = "15px solid #93D2B5";
                            } else {
                                label_type.style.borderTop = "15px solid #E75C8D";
                                label_type.style.borderBottom = "0px solid transparent";
                            }
                            label_weight_value.innerHTML = weight_value_diff + "&nbsp;kg";
                        }
                    }
                }
            });
    });
}

// Handle incoming messages. Called when:
// - a message is received while the app has focus
// - the user clicks on an app notification created by a service worker
//   `messaging.onBackgroundMessage` handler.
// signin画面でcallされるのを確認
// onMessage(messaging, (payload) => {
//     console.log('onMessage messaging called. ', payload);
// });