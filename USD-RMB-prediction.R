plotForecastErrors <- function(forecasterrors)
  {
     # make a histogram of the forecast errors:
     mybinsize <- IQR(forecasterrors)/4
     mysd   <- sd(forecasterrors)
     mymin  <- min(forecasterrors) - mysd*5
     mymax  <- max(forecasterrors) + mysd*3
     # generate normally distributed data with mean 0 and standard deviation mysd
     mynorm <- rnorm(10000, mean=0, sd=mysd)
     mymin2 <- min(mynorm)
     mymax2 <- max(mynorm)
     if (mymin2 < mymin) { mymin <- mymin2 }
     if (mymax2 > mymax) { mymax <- mymax2 }
     # make a red histogram of the forecast errors, with the normally distributed data overlaid:
     mybins <- seq(mymin, mymax, mybinsize)
     hist(forecasterrors, col="red", freq=FALSE, breaks=mybins)
     # freq=FALSE ensures the area under the histogram = 1
     # generate normally distributed data with mean 0 and standard deviation mysd
     myhist <- hist(mynorm, plot=FALSE, breaks=mybins)
     # plot the normal curve as a blue line on top of the histogram of forecast errors:
     points(myhist$mids, myhist$density, type="l", col="blue", lwd=2)
  }

exchangeRate <- read.csv("./USDtoRMB.csv",header=TRUE, sep=",")
exchangeRateTS <- ts(exchangeRate$Mid, frequency=365, start=c(2012,1))
plot.ts(exchangeRateTS)

install.pakcages("TTR")
library("TTR")

exchangeRateSMA5 <- SMA(exchangeRateTS, n=5)
plot.ts(exchangeRateSMA5)

exchangeRateSMA10 <- SMA(exchangeRateTS, n=10)
plot.ts(exchangeRateSMA10)

exchangeRateSMA30 <- SMA(exchangeRateTS, n=30)
plot.ts(exchangeRateSMA30)

exchangeRateDecomposed <- decompose(exchangeRateTS)
plot(exchangeRateDecomposed)

exchangeRateForecast <- HoltWinters(exchangeRageTS, gamma=FALSE)
plot(exchangeRateForecaset)

install.packages("forecast")
library(forecast)

exchangeRateForecast2 <- forecast.HoltWinters(exchangeRateForecast, h=180)

exchangeRateArima <- arima(exchangeRateTS, order=c(0,1,0))
exchangeRateArimaFore <- forecast.Arima(exchangeRateArima, h=180)
plot.forecast(exchangeRateArimaFore)
exchangeRatediff <- diff(ts, differences=1)
plot.ts(exchangeRatediff)


install.packages("fGarch")
tsgarchFit<-garchFit(formula = ~garch(1,1),data=ts)
predict(tsgarchFit, n.ahead=180,trace=TRUE,plot=TRUE)
