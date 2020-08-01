package com.prueba.myapplication;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;

import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.Volley;
//import com.prueba.myapplication.models.Contacto;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class MainActivity extends AppCompatActivity {

    //ListView contactosList;
    //ContactoAdaptador contactoAdaptador;
    RequestQueue requestQueue;
    String url = "http://192.168.0.3:4000/listar";
    TextView historial;
    //Button boton;
    EditText dominio;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        historial = findViewById(R.id.historial_endpoint);
        dominio = findViewById(R.id.edit_text);
        requestQueue = Volley.newRequestQueue(this);

    }

    public void consultaHistorial(View view){
        System.out.println("________________________________________________________________________________");
        JsonObjectRequest request = new JsonObjectRequest(Request.Method.GET, url, null, new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                            try {
                                JSONArray list = response.getJSONArray("his");
                                String respuesta = "Dominios \n\n";
                                System.out.println(respuesta);
                                for (int i=0; i < list.length(); i++) {
                                    JSONObject o = list.getJSONObject(i);
                                    respuesta += o.getString("domain") + "\n";
                                }
                                historial.setText(respuesta);
                            } catch (JSONException e) {
                                e.printStackTrace();
                            }
                    }
                }, new Response.ErrorListener() {

                    @Override
                    public void onErrorResponse(VolleyError error) {

                    }
                });
        requestQueue.add(request);
    }

    public void consultaDominio(View view){
        url = "http://192.168.0.3:4000/whois/" + dominio.getText().toString();
        JsonObjectRequest request = new JsonObjectRequest(Request.Method.GET, url, null, new Response.Listener<JSONObject>() {
            @Override
            public void onResponse(JSONObject response) {
                try {
                    JSONArray list = response.getJSONArray("List");
                    String respuesta = "Servers \n\n";
                    System.out.println(respuesta);
                    for (int i=0; i < list.length(); i++) {
                        JSONObject o = list.getJSONObject(i);
                        respuesta += o.getString("servers") + "\n";
                    }
                    historial.setText(respuesta);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        }, new Response.ErrorListener() {

            @Override
            public void onErrorResponse(VolleyError error) {

            }
        });
        requestQueue.add(request);
    }

}
