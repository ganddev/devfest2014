package com.ganddev.devfest2014.utils;

import android.content.Context;

import com.android.volley.AuthFailureError;
import com.android.volley.Response;
import com.android.volley.toolbox.JsonObjectRequest;

import org.json.JSONObject;

import java.util.HashMap;
import java.util.Map;

/**
 * Created by bjornahlfeld on 23.11.14.
 */
public class CustomRequest extends JsonObjectRequest {

    private Context mContext;
    private String deviceToken;

    public CustomRequest(int method, String url, JSONObject jsonRequest,
                       Response.Listener<JSONObject> listener,
                       Response.ErrorListener errorListener, Context ctx, String deviceToken) {
        super(method, url, jsonRequest, listener, errorListener);
        mContext = ctx;
        this.deviceToken = deviceToken;
    }

    public CustomRequest(String url, JSONObject jsonRequest,
                       Response.Listener<JSONObject> listener,
                       Response.ErrorListener errorListener, Context ctx) {
        super(url, jsonRequest, listener, errorListener);
        mContext = ctx;
    }

    @Override
    public Map<String, String> getHeaders() throws AuthFailureError {
        return createBasicAuthHeader();
    }

    Map<String, String> createBasicAuthHeader() {
        Map<String, String> headerMap = new HashMap<String, String>();
        headerMap.put("deviceToken", deviceToken);
        headerMap.put("Content-Type", "application/json");
        return headerMap;
    }
}
