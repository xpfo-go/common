# GIS

## 一 介绍

### 地理信息系统，这里提供各个地图经纬度之间转换

## 二 使用

#### WGS84坐标系：即地球坐标系，国际上通用的坐标系。

#### GCJ02坐标系：即火星坐标系，WGS84坐标系经加密后的坐标系。Google Maps，高德在用。

#### BD09坐标系：即百度坐标系，GCJ02坐标系经加密后的坐标系。

``` go
// BD09toGCJ02 百度坐标系->火星坐标系
func BD09toGCJ02(lng, lat float64) (float64, float64) {}

// GCJ02toBD09 火星坐标系->百度坐标系
func GCJ02toBD09(lng, lat float64) (float64, float64) {}

// WGS84toGCJ02 WGS84坐标系->火星坐标系
func WGS84toGCJ02(lng, lat float64) (float64, float64) {}

// GCJ02toWGS84 火星坐标系->WGS84坐标系
func GCJ02toWGS84(lng, lat float64) (float64, float64) {}

// BD09toWGS84 百度坐标系->WGS84坐标系
func BD09toWGS84(lng, lat float64) (float64, float64) {}

// WGS84toBD09 WGS84坐标系->百度坐标系
func WGS84toBD09(lng, lat float64) (float64, float64) {}
```
