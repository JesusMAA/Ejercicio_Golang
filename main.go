package main

/*
  **importar packages
  *package para usar la codificacion y decodificacion de json
  *log -manejo de errores en el servidor
  *net/http para manejar las peticiones del servidor
  *
*/
import (
  "encoding/json"
  "log"
  "github.com/jpiontek/go-ip-api"//api para saver la ip
  "fmt"
  "strings"
  "github.com/fasthttp/router"
  "github.com/valyala/fasthttp"
  "github.com/keltia/ssllabs"
  "database/sql"
  _ "github.com/lib/pq"
  "time"

)

type Person struct {
  ID string `jason:"id,omitempy"`
  FirstName string `jason:"firtname,omitempy"`
  LastName string `jason:"lastname,omitempy"`
  Address *Address `jason:"address,omitempy"`
}

type Address struct{
  City string `jason:"city,omitempy"`
  State string `jason:"state,omitempy"`
}

type EndPoints struct{
  List []ServersList `jason:"list,omitempy"`
  Info *InfoServer  `jason:"infoServer, omitempy"`
}

type ServersList struct{
  Servers []Servers `json:"servers,omitempy"`
}

type Servers struct{
  Address string `jason:"address,omitempy"`
  Ssl string `jason:"ssl,omitempy"`
  Country string `jason:"country,omitempy"`
  Owner string `jason:"owner,omitempy"`
}


type InfoServer struct {
	Servers_changed     string    `json:"server_changed,omitempy"`
	Ssl_grade          string    `json:"ssl_grade,omitempy"`
	Previous_ssl_grade string    `json:"previous_ssl_grade,omitempy"`
	Logo               string    `json:"logo,omitempy"`
	Title              string    `json:"title,omitempy"`
	Is_down            string    `json:"is_down,omitempy"`
}

type History struct{
  His []Domain  `json:"his,omitempy"`
}

type Domain struct{
  Domain string `json:"domain,omitempy"`
}

//arreglo de personas simulando persistencia
var people []Person
var serversArr []Servers
var serverListTemp []ServersList


//request -> para cuando el navegador envie la informacion
//viene desde el cliente
func GetPeopleEndPoint(request *fasthttp.RequestCtx){
  //el modulo que se importo
  json.NewEncoder(request).Encode(people)
}

func GetWhoisEndPoint(request *fasthttp.RequestCtx){
  serversArr = nil
  serverListTemp = nil
  //el modulo que se importo
    params := request.UserValue("id").(string)


    client, _ := ssllabs.NewClient()


    opts := map[string]string{"all": ""}
    lr, _ := client.Analyze(params, false, []map[string]string{opts}...)

    temp := goip.NewClient()

    for i:= 0; i<len(lr.Endpoints); i++{
      result, _ := temp.GetLocationForIp(lr.Endpoints[i].IPAddress)

      serversArr = append(serversArr, Servers{Address: lr.Endpoints[i].IPAddress,
         Ssl:lr.Endpoints[i].Grade, Country:result.Country, Owner:result.Isp})
    }
    aux24 := ServersList{Servers: serversArr}
    //recorrer aux24.Servers porque es un arreglo
    if aux24.Servers == nil{
      fmt.Print("hola")

      insertDb(request)
      serversArr = existDb(serversArr, request)
      grade_low := gradeLow(serversArr)
      fmt.Println("____________________UNO_____________________________--")
      server_changed, ssl_grade := serverChanged(serversArr, params)
      status := ""

      if len(lr.Endpoints)>0{
        status = lr.Endpoints[0].StatusMessage
      }


      serverListTemp = append(serverListTemp, ServersList{Servers: serversArr})
      GetInfoServer(params, request, serverListTemp, grade_low, server_changed, ssl_grade, status)
      //fmt.Println(serversArr)
    }else
    {

      serverListTemp = append(serverListTemp, ServersList{Servers: serversArr})

      grade_low := gradeLow(serversArr)

      server_changed, ssl_grade := serverChanged(serversArr, params)
      status := ""
      if len(lr.Endpoints)>0{
        status = lr.Endpoints[0].StatusMessage
      }


      //setServers()
      insertDb(request)
      insertServersDb(serversArr, params)
      fmt.Println("________________DOS_________________________________--")
      GetInfoServer(params, request, serverListTemp, grade_low, server_changed, ssl_grade, status)
  }
}

func GetInfoServer(url string, request *fasthttp.RequestCtx, arr []ServersList, grade_low string, server_changed string, ssl_grade string, status string){
  var dst []byte
  _, resp,_ := fasthttp.Get(dst,"https://i.olsh.me/icons?url="+url)
  posicion := strings.Index(string(resp), "<td><a href=")
  posicion2 := 100

  var linkIcon string = ""

  for i := posicion+13;  i < posicion2 || string(resp[i])!= ">"; i++ {
    linkIcon +=string(resp[i])
  }
  linkIcon = strings.TrimRight(linkIcon, "'")
fmt.Println("________________TRES_________________________________--")
  //Obtener title
  var dst2 []byte
  _, resp2,_ := fasthttp.Get(dst2,"https://"+url)
  posicion2 = strings.Index(string(resp2), "</title>")
  posicion = strings.Index(string(resp2), "<title")

  pageTitle := []byte(string(resp2[posicion:posicion2]))

  posicion = strings.Index(string(pageTitle), ">")
  posicion += 1
  tam := len(string(pageTitle))

  titleAux := []byte(string(pageTitle[posicion:tam]))
  title := string(titleAux)

  

   info := &InfoServer{Servers_changed: server_changed, Ssl_grade:grade_low,
  Previous_ssl_grade:ssl_grade, Logo:linkIcon, Title:title, Is_down:status}


  aux :=  EndPoints{List: arr, Info: info}

  json.NewEncoder(request).Encode(aux)
}

func gradeLow(arr []Servers)string{
  grade := []string{"A+", "A", "B", "C", "D", "E", "F"}
	gradeByte := strings.Join(grade, ",")


  var cont int = 0
  aux := 0

  for i := 0;  i< len(arr) && aux != -1; i++ {
    aux = strings.Index(gradeByte,arr[i].Ssl+",")
    cont ++
  }

  return grade[cont]
}

func serverChanged(arr []Servers, domain string) (string, string){

  db, err := sql.Open("postgres", "postgres://root@localhost:26257/servers?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT domain, time FROM domains WHERE domain='"+ domain +"'")

    if err != nil {
    		log.Fatal(err)
    }
    defer rows.Close()

    var addressDb, timeDb string
    for rows.Next() {
        if err := rows.Scan(&addressDb, &timeDb); err != nil {
            log.Fatal(err)
        }
    }

    aux := strings.Index(timeDb, "T")
    fmt.Println(timeDb)
    timeDb1 := timeDb[aux+1] + timeDb[aux+2]

    aux2 := time.Now().String()
    posicion := strings.Index(aux2, " -")
    titleAux := string(aux2[0:posicion])
    aux = strings.Index(titleAux, " ")
    timeDb2 := timeDb[aux+1] + timeDb[aux+2]

    //comparacion de servers
    if timeDb1 != timeDb2{
      rows, err = db.Query("SELECT domain FROM server WHERE domain='"+ domain +"'")
      if err != nil {
      		log.Fatal(err)
      }
      defer rows.Close()

      i :=0

      for rows.Next() {
          var addressDb string

          if err := rows.Scan(&addressDb); err != nil {
              log.Fatal(err)
          }
          //CONTINUAR AQUI
          if arr[i].Address != addressDb && i<len(arr){
            return "true", arr[i].Ssl
          }
          i++
      }
    }

      return "false", "nil"
}

//llenar servers struct
func existDb(arr []Servers, request *fasthttp.RequestCtx) []Servers{
  domain := request.UserValue("id").(string)
  db, err := sql.Open("postgres", "postgres://root@localhost:26257/servers?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT addres, ssl_grade, country, owner FROM server WHERE domain='"+ domain +"'")
    if err != nil {
    		log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var addressDb, ssl_gradeDb, countryDb, ownerDb string

        if err := rows.Scan(&addressDb, &ssl_gradeDb, &countryDb, &ownerDb); err != nil {
            log.Fatal(err)
        }
        //CONTINUAR AQUI
        arr = append(arr, Servers{Address: addressDb,
           Ssl:ssl_gradeDb, Country:countryDb, Owner:ownerDb})
        //fmt.Println( domainDb, timeDb)
    }
    return arr
}

func insertServersDb(arr []Servers, domain string) {
  db, err := sql.Open("postgres", "postgres://root@localhost:26257/servers?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }
    defer db.Close()

    _, err = db.Query("SELECT domain FROM server WHERE domain="+ domain +"")
    //fmt.Println("SELECT domain FROM server WHERE domain='"+ domain +"'")
    if err != nil{
      for i:= 0; i<len(arr); i++{
        if _, err := db.Exec(
          "INSERT INTO server (domain, addres, ssl_grade, country, owner) VALUES ('" + domain + "','"+ arr[i].Address + "','" + arr[i].Ssl + "','" + arr[i].Country + "','" + arr[i].Owner + "')");
          err != nil {
      		fmt.Println("esto no es un error")
      	}
      }
    }

}

//guarda los dominios
func insertDb(request *fasthttp.RequestCtx) {

  domain := request.UserValue("id").(string)
  db, err := sql.Open("postgres", "postgres://root@localhost:26257/servers?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }
    defer db.Close()
    aux := time.Now().String()
    posicion := strings.Index(aux, " -")
    titleAux := string(aux[0:posicion])
    _, err = db.Query("SELECT domain FROM domains WHERE domain='"+ domain +"'")


    fmt.Println(titleAux)
    fmt.Println(err)

    if err != nil {
      if _, err := db.Exec(
    		"INSERT INTO domains (domain, time) VALUES ('" + domain + "','" + titleAux + "')");
        err != nil {
    		fmt.Println("Hola")
    	}
    }

    //defer rows.Close()


}

func GetLista(request *fasthttp.RequestCtx)  {

  var arr []Domain
  //var arr2 []History

  db, err := sql.Open("postgres", "postgres://root@localhost:26257/servers?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT domain FROM domains")
    if err != nil {
    		log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var addressDb string

        if err := rows.Scan(&addressDb); err != nil {
            log.Fatal(err)
        }
        //CONTINUAR AQUI
        arr = append(arr, Domain{Domain:addressDb})
        //fmt.Println( domainDb, timeDb)
    }
    arr2 := History{His:arr}
    json.NewEncoder(request).Encode(arr2)
}

func main()  {

  rout := router.New()
  //fmt.Println(time.Now())

  //GetData()



  rout.GET("/whois/{id}", GetWhoisEndPoint)
  rout.GET("/listar", GetLista)




  //se usa "log" en caso de que ocurra un error
  log.Fatal(fasthttp.ListenAndServe(":4000", rout.Handler))

}
