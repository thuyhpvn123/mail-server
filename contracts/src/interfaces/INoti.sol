// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity 0.8.19;

struct NotiParams{
    string title;
    string body;
}
enum PlatformEnum{
    ANDROID,
    IOS,
    WEB
}
interface INoti{

function AddNoti(
    NotiParams memory params,
    address _to
  ) external returns (bool);

}
