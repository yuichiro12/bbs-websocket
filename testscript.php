<?php

$sock = "unix:///Users/yuichiro/Sites/bbs/tmp/bbs-uds.sock";
$fp = fsockopen($sock);
$json = [
	"message" => "movさんが重症です。",
    "url" => "https://bbs-localhost/image/users/5a253d71b6b7c858171528.png",
    "icon" => "https://bbs-localhost/image/users/5a253d71b6b7c858171528.png",
    "ids" => [6],
];

fwrite($fp, json_encode($json));
fclose($fp);
