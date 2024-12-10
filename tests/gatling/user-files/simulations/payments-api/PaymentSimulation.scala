import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._

class Payment extends Simulation {

  val httpProtocol = http
    .baseUrl("http://payment-transaction-rest:8080")
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")


  val paymentExecute = 
    feed(tsv("transactions.tsv").circular())
    .exec(session => {
      val payload = session("payload").as[String]
      println(s"Payload JSON: $payload")
      session
    })
    .exec(http("Requisição para /payments")
      .post("/payment")
      .body(StringBody("#{payload}"))
      .header("Content-Type", "application/json")
      .check(status.in(200))
    )

  val testPaymentExecute = scenario("Test Payments").exec(paymentExecute)


  /*
    Common TPS:
     9000k in 5min -  30 TPS
    15000k in 5min -  50 TPS
    30000k in 5min - 100 TPS
  */
  setUp(
    testPaymentExecute.inject(rampUsers(15000).during(301.seconds)) // 7500k in 5min - 25 TPS
  ).protocols(httpProtocol)
}
