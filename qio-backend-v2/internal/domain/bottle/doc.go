// Package bottle 承载漂流瓶域：投放、捞取、扔回。
//
// 通过漂流瓶发起好友申请不在本域——那是 bottle 与 friend 两个域的协作，
// v1 把它放在 BottleServiceImpl 中，v2 归 internal/app。
//
// 瓶子图片的渲染同样不在本域，属于展示层关注点，渲染方案待定。
package bottle
