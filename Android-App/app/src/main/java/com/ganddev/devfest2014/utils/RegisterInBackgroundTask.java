package com.ganddev.devfest2014.utils;

import android.content.Context;
import android.os.AsyncTask;
import android.util.Log;

import com.ganddev.devfest2014.R;
import com.ganddev.devfest2014.api.ApiService;
import com.google.android.gms.gcm.GoogleCloudMessaging;

import java.io.IOException;

/**
 * Created by bjornahlfeld on 22.11.14.
 */
public class RegisterInBackgroundTask extends AsyncTask<Void, Void, String > {

    private static final String TAG = RegisterInBackgroundTask.class.getSimpleName();
    private Context mContext;
    private GoogleCloudMessaging gcm;

    public RegisterInBackgroundTask(Context ctx)
    {
        mContext = ctx;
        gcm = GoogleCloudMessaging.getInstance(ctx);
    }

    @Override
    protected String doInBackground(Void... params) {
        String msg = "";
        try {
            if (gcm == null) {
                gcm = GoogleCloudMessaging.getInstance(mContext);
            }
            final String regid = gcm.register(mContext.getString(R.string.gcm_sender_id));
            msg = "Device registered, registration ID=" + regid;

            // You should send the registration ID to your server over HTTP,
            // so it can use GCM/HTTP or CCS to send messages to your app.
            // The request to your server should be authenticated if your app
            // is using accounts.
            sendRegistrationIdToBackend(regid);

            // For this demo: we don't need to send it because the device
            // will send upstream messages to a server that echo back the
            // message using the 'from' address in the message.

            // Persist the regID - no need to register again.
            SharedPreferencesUtils.storeRegistrationId(mContext, regid);
        } catch (IOException ex) {
            msg = "Error :" + ex.getMessage();
            // If there is an error, don't just keep trying to register.
            // Require the user to click a button again, or perform
            // exponential back-off.
        }
        return msg;
    }

    private void sendRegistrationIdToBackend(final String regid) {
        ApiService.postGCMRegId(regid);
    }

    @Override
    protected void onPostExecute(String msg) {
        Log.d(TAG, "REG_ID: " +msg);
        ApiService.postGCMRegId(msg);
    }
}
